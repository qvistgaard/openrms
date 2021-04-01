package oxigen

import (
	"context"
	"errors"
	"fmt"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"openrms/ipc"
	"openrms/ipc/commands"
	"time"
)

type Oxigen struct {
	state    byte
	serial   io.ReadWriteCloser
	settings Settings
	version  string
	running  bool
}

type Settings struct {
	maxSpeed byte
	pitLane  PitLane
}

type PitLane struct {
	lapCounting byte
	lapTrigger  byte
}

func Connect(serial io.ReadWriteCloser) (*Oxigen, error) {
	var err error
	o := new(Oxigen)
	o.serial = serial
	o.running = true

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		o.serial.Close()
		return nil, err
	}

	versionResponse := make([]byte, 13)
	time.Sleep(10 * time.Millisecond)
	_, err = o.serial.Read(versionResponse)
	v, _ := version.NewVersion(fmt.Sprintf("%d.%d", versionResponse[0], versionResponse[1]))
	constraint, _ := version.NewConstraint(">= 3.10")

	if !constraint.Check(v) {
		return nil, errors.New(fmt.Sprintf("Unsupported dongle version: %s", v))
	}
	o.version = v.Original()
	log.Printf("Connected to oxigen dongle. Dongle version: %s", v)

	return o, nil
}

func (oxigen *Oxigen) Close() error {
	return oxigen.serial.Close()
}

func (oxigen *Oxigen) EventLoop(input queue.Queue, output queue.Queue) error {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	defer oxigen.serial.Close()
	var err error

	for {
		var command *ipc.Command
		cmd, _ := input.DequeueOrWaitForNextElementContext(ctx)
		if cmd == nil {
			command = ipc.NewEmptyCommand()
		} else {
			command = cmd.(*ipc.Command)
		}
		b := oxigen.command(*command)
		_, err = oxigen.serial.Write(b)
		if err != nil {
			break
		}
		// log.Printf("S> %d %s", len(b), hex.Dump(b))
		for {
			time.Sleep(10 * time.Millisecond)
			buffer := make([]byte, 13)
			_, err := oxigen.serial.Read(buffer)
			// log.Printf("S> %d %s", len(buffer), hex.Dump(buffer))

			event := oxigen.event(buffer)
			if event.Id > 0 {
				output.Enqueue(event)
			} else {
				break
			}
			if err != nil {
				log.Print(err)
				err = nil
				break
			}
		}
		if err != nil {
			break
		}
		if !oxigen.running {
			break
		}
	}
	log.Printf("error: %s", err)
	return err
}

func (oxigen *Oxigen) event(b []byte) ipc.Event {
	return ipc.Event{
		Id: b[1],
		Controller: ipc.Controller{
			BatteryWarning: 0x04&b[0] == 0x04,
			Link:           0x02&b[0] == 0x02,
			TrackCall:      0x08&b[0] == 0x08,
			ArrowUp:        0x20&b[0] == 0x20,
			ArrowDown:      0x40&b[0] == 0x40,
			// version: ,
		},
		Car: ipc.Car{
			Reset: 0x01&b[0] == 0x01,
			InPit: 0x40&b[8] == 0x40,
			// version:
		},
		LapTime:      time.Duration((uint(b[2]) * 256) + uint(b[3])),
		LapNumber:    (uint16(b[5]) * 256) + uint16(b[6]),
		TriggerValue: 0x74 & b[7],
		Ontrack:      0x80&b[7] == 0x80,
	}
}

func (oxigen *Oxigen) command(c ipc.Command) []byte {
	var cmd byte = 0x00
	var parameter byte = 0x00
	var controller byte = 0x00

	switch c.CommandType().(type) {
	case *commands.MaxSpeed:
		cmd = 0x02 // TODO: Add global command support
		controller = c.Driver()
		parameter = c.Value()[0]
	}

	return []byte{
		oxigen.state | oxigen.settings.pitLane.lapTrigger | oxigen.settings.pitLane.lapCounting,
		oxigen.settings.maxSpeed,
		controller,
		cmd,
		parameter,
		0x00, // unused
		0x00, // unused
		0x00, // Racetimer ? ? TODO: Figure this out
		0x00, // Racetimer ? ?
		0x00, // Racetimer ? ?
	}
}

func (oxigen *Oxigen) MaxSpeed(speed uint8) {
	oxigen.settings.maxSpeed = speed
}

func (oxigen *Oxigen) Start() {
	oxigen.state = 0x03
}

func (oxigen *Oxigen) PitLaneLapCount(enabled bool, entry bool) {
	if !enabled {
		oxigen.settings.pitLane.lapCounting = 0x20
		oxigen.settings.pitLane.lapTrigger = 0x00
	} else {
		oxigen.settings.pitLane.lapCounting = 0x00
		if entry {
			oxigen.settings.pitLane.lapTrigger = 0x00
		} else {
			oxigen.settings.pitLane.lapTrigger = 0x40
		}
	}
}

func (oxigen *Oxigen) Stop() {
	oxigen.state = 0x01
}

func (oxigen *Oxigen) Pause() {
	oxigen.state = 0x04
}

func (oxigen *Oxigen) Flag(lc bool) {
	if lc {
		oxigen.state = 0x05
	} else {
		oxigen.state = 0x15
	}
}

func (oxigen *Oxigen) Shutdown() {
	oxigen.running = false
}
