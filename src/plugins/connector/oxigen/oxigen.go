package oxigen

import (
	"../../../ipc"
	"../../../ipc/commands"
	"context"
	"errors"
	"fmt"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"time"
	// "../../../src/connector"
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
	o.stop()
	o.serial = serial
	o.running = true

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		o.serial.Close()
		return nil, err
	}

	versionResponse := make([]byte, 13)
	_, err = o.serial.Read(versionResponse)
	v, _ := version.NewVersion(fmt.Sprintf("%d.%d", versionResponse[0], versionResponse[1]))
	constraint, _ := version.NewConstraint(">= 3.0")

	if !constraint.Check(v) {
		return nil, errors.New(fmt.Sprintf("Unsupported dongle version: %s", v))
	}
	o.version = v.Original()
	log.Printf("Connected to oxigen dongle. Dongle version: %s", v)

	return o, nil
}

func (oxigen *Oxigen) Closer() error {
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
		fmt.Printf("COMMAND [%T] %+v\n", command, command)

		b := oxigen.message(*command)
		_, err = oxigen.serial.Write(b)
		if err != nil {
			break
		}
		log.Printf("S> %x", b)
		for {
			buffer := make([]byte, 13)
			len, err := oxigen.serial.Read(buffer)

			if len > 0 {
				log.Printf("R> %x", buffer)

				qerr := output.Enqueue(buffer)
				if qerr != nil {
					break
				}
			}

			// os.Exit(1)
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

func (oxigen *Oxigen) message(c ipc.Command) []byte {
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

func (oxigen *Oxigen) maxSpeed(speed uint8) bool {
	oxigen.settings.maxSpeed = speed
	return true
}

func (oxigen *Oxigen) start() bool {
	oxigen.state = 0x03
	return true
}

func (oxigen *Oxigen) pitLaneLapCount(enabled bool, entry bool) bool {
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
	return true
}

func (oxigen *Oxigen) stop() bool {
	oxigen.state = 0x01
	return true
}

func (oxigen *Oxigen) pause() bool {
	oxigen.state = 0x04
	return true
}

func (oxigen *Oxigen) flag(lc bool) bool {
	if lc {
		oxigen.state = 0x05
	} else {
		oxigen.state = 0x15
	}
	return true
}

func (oxigen *Oxigen) shutdown() {
	oxigen.running = false
}
