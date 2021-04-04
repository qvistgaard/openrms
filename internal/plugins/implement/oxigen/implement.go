package oxigen

import (
	"context"
	"errors"
	"fmt"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/hashicorp/go-version"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

type Oxigen struct {
	state    byte
	serial   io.ReadWriteCloser
	settings *Settings
	version  string
	running  bool
	commands queue.Queue
	events   queue.Queue
}

func CreateUSBConnection(device string) (*serial.Port, error) {
	c := &serial.Config{
		Name:        device,
		Baud:        19200,
		ReadTimeout: time.Millisecond * 100,
	}
	return serial.OpenPort(c)
}

func CreateImplement(serial io.ReadWriteCloser) (*Oxigen, error) {
	var err error
	o := new(Oxigen)
	o.serial = serial
	o.running = true
	o.commands = queue.NewFIFO()
	o.events = queue.NewFIFO()
	o.settings = newSettings()

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

func (o *Oxigen) EventLoop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	defer o.serial.Close()
	var err error
	timer := []byte{0x00, 0x00, 0x00}
	for {
		var command *Command
		cmd, _ := o.commands.DequeueOrWaitForNextElementContext(ctx)

		if cmd == nil {
			command = newEmptyCommand(map[string]state.StateInterface{}, o.state, o.settings)
		} else {
			command = cmd.(*Command)
		}

		b := o.command(command, timer)
		_, err = o.serial.Write(b)
		if err != nil {
			break
		}
		// log.Printf("S> %d %s", len(b), hex.Dump(b))
		for {
			time.Sleep(10 * time.Millisecond)
			buffer := make([]byte, 13)
			_, err := o.serial.Read(buffer)
			// log.Printf("S> %d %s", len(buffer), hex.Dump(buffer))
			timer = buffer[7:9]

			event := o.event(buffer)
			if event.Id > 0 {
				o.events.Enqueue(event)
			} else {
				break
			}
			if err != nil {
				log.Print(err)
				err = nil
				break
			}
		}
		if !o.running {
			break
		}
	}
	log.Printf("error: %s", err)
	return err
}

func (o *Oxigen) WaitForEvent() (implement.Event, error) {
	element, err := o.events.DequeueOrWaitForNextElement()
	return element.(implement.Event), err
}

func (o *Oxigen) SendCommand(c implement.Command) error {
	if len(c.Changes.Car) > 0 {
		for k, v := range c.Changes.Car {
			ec := newEmptyCommand(c.Changes.Race, o.state, o.settings)
			if ec.carCommand(c.Id, k, v) {
				err := o.commands.Enqueue(ec)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return o.commands.Enqueue(newEmptyCommand(c.Changes.Race, o.state, o.settings))
	}
}

func (o *Oxigen) event(b []byte) implement.Event {
	return implement.Event{
		Id: b[1],
		Controller: implement.Controller{
			BatteryWarning: 0x04&b[0] == 0x04,
			Link:           0x02&b[0] == 0x02,
			TrackCall:      0x08&b[0] == 0x08,
			ArrowUp:        0x20&b[0] == 0x20,
			ArrowDown:      0x40&b[0] == 0x40,
			// version: ,
		},
		Car: implement.Car{
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

func (o *Oxigen) command(c *Command, timer []byte) []byte {
	var cmd byte = 0x00
	var parameter byte = 0x00
	var controller byte = 0x00

	if c.car != nil {
		cmd = c.car.command
		parameter = c.car.value
		controller = c.car.id
	}

	return []byte{
		o.state | o.settings.pitLane.lapTrigger | o.settings.pitLane.lapCounting,
		o.settings.maxSpeed,
		controller,
		cmd,
		parameter,
		0x00,     // unused
		0x00,     // unused
		timer[0], // Racetimer
		timer[1], // Racetimer
		timer[2], // Racetimer
	}
}
