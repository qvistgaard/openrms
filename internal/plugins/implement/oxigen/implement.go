package oxigen

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"io"
	"time"
)

type Oxigen struct {
	state      byte
	serial     io.ReadWriteCloser
	settings   *Settings
	version    string
	running    bool
	commands   chan *Command
	events     chan implement.Event
	bufferSize int
	timer      []byte
	buffer     *bufio.Reader
}

func CreateUSBConnection(device string) (*serial.Port, error) {
	c := &serial.Config{
		Name:   device,
		Baud:   115200,
		Parity: serial.ParityNone,
	}
	return serial.OpenPort(c)
}

func CreateImplement(serial io.ReadWriteCloser) (*Oxigen, error) {
	var err error
	o := new(Oxigen)
	o.serial = serial
	o.buffer = bufio.NewReaderSize(o.serial, 1024)
	o.bufferSize = 1024
	o.commands = make(chan *Command, 1024)
	o.events = make(chan implement.Event, 1024)
	o.settings = newSettings()

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		o.serial.Close()
		return nil, err
	}

	versionResponse := make([]byte, 2)
	_, err = io.ReadFull(o.buffer, versionResponse)
	v, _ := version.NewVersion(fmt.Sprintf("%d.%d", versionResponse[0], versionResponse[1]))
	constraint, _ := version.NewConstraint(">= 3.10")

	log.WithField("message", fmt.Sprintf("%x", versionResponse)).
		WithField("bytes", len(versionResponse)).
		Trace("received message from oxygen dongle")

	if !constraint.Check(v) {
		return nil, errors.New(fmt.Sprintf("Unsupported dongle version: %s", v))
	}
	o.version = v.Original()
	log.WithField("version", v).Infof("Connected to oxigen dongle. Dongle version: %s", v)

	return o, nil
}

func (o *Oxigen) EventLoop() error {
	defer func() {
		o.serial.Close()
		panic("oxigen: keep-alive routine failed.")
	}()

	o.timer = []byte{0x00, 0x00, 0x00}

	// For some reason unkown, the dongle sends more data back when requesting
	// the version number, therefore we discard what ever is buffered before
	// we start the real work
	o.buffer.Discard(o.buffer.Buffered())

	// Keep-alive routine, to keep oxigen dongle sending data back
	go o.keepAlive()
	go o.sendCommand()

	for {
		buffer := make([]byte, 13)
		read, err := io.ReadFull(o.buffer, buffer)

		log.WithField("message", fmt.Sprintf("%x", buffer)).
			WithField("bytes", read).
			Trace("received message from oxygen dongle")
		if err == nil {
			if read == 13 {
				o.timer = buffer[7:10]
				o.events <- o.event(buffer)
			}
		} else {
			log.Error(err)
		}
	}
}

func (o *Oxigen) sendCommand() {
	defer func() {
		panic("oxigen: keep-alive routine failed.")
	}()
	for {
		select {
		case cmd := <-o.commands:
			if float32(len(o.commands)) > (float32(o.bufferSize) * 0.9) {
				log.WithFields(map[string]interface{}{
					"bufferSize": o.bufferSize,
					"size":       len(o.commands),
				}).Warn("too many commands on command buffer")
			}

			b := o.command(cmd, o.timer)
			o.settings = &cmd.settings
			o.state = cmd.state
			l, err := o.serial.Write(b)
			if err != nil {
				panic(err)
			}
			if cmd.car != nil {
				log.WithFields(map[string]interface{}{
					"message": fmt.Sprintf("%x", b),
					"size":    fmt.Sprintf("%d", l),
					"decode":  hex.Dump(b),
				}).Debug("send message to oxygen dongle")
			} else {
				log.WithFields(map[string]interface{}{
					"message": fmt.Sprintf("%x", b),
					"size":    fmt.Sprintf("%d", l),
				}).Trace("send message to oxygen dongle")

			}
		}
	}
}

func (o *Oxigen) keepAlive() {
	defer func() {
		panic("oxigen: keep-alive routine failed.")
	}()
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			if len(o.commands) == 0 {
				o.commands <- newEmptyCommand(state.CourseState{}, o.state, o.settings)
				log.Trace("oxigen: sent keep-alive")
			}
		}
	}
}

func (o *Oxigen) EventChannel() <-chan implement.Event {
	return o.events
}

func (o *Oxigen) SendRaceState(r state.CourseState) error {
	o.commands <- newEmptyCommand(r, o.state, o.settings)
	return nil
}

func (o *Oxigen) SendCarState(c state.CarState) error {
	if len(c.Changes) > 0 {
		for _, v := range c.Changes {
			ec := newEmptyCommand(state.CourseState{}, o.state, o.settings)
			if ec.carCommand(uint8(c.Car), v.Name, v.Value) {
				o.commands <- ec
			}
		}
	}
	return nil
}

func (o *Oxigen) ResendCarState(c *state.Car) {
	resendStates := []string{
		state.CarMaxSpeed, state.CarMaxBreaking, state.CarMinSpeed, state.CarPitLaneSpeed,
	}
	for _, n := range resendStates {
		ec := newEmptyCommand(state.CourseState{}, o.state, o.settings)
		if ec.carCommand(uint8(c.Id()), n, c.Get(n)) {
			o.commands <- ec
		}
	}
}

func (o *Oxigen) event(b []byte) implement.Event {
	rt := (uint(b[9]) * 16777216) + (uint(b[10]) * 65536) + (uint(b[11]) * 256) + uint(b[12]) - uint(b[4])
	rtd := time.Duration(rt*10) * time.Millisecond
	e := implement.Event{
		Id: state.CarId(b[1]),
		Controller: implement.Controller{
			BatteryWarning: 0x04&b[0] == 0x04,
			Link:           0x02&b[0] == 0x02,
			TrackCall:      0x08&b[0] == 0x08,
			ArrowUp:        0x20&b[0] == 0x20,
			ArrowDown:      0x40&b[0] == 0x40,
		},
		Car: implement.Car{
			Reset: 0x01&b[0] == 0x01,
			InPit: 0x40&b[8] == 0x40,
		},
		Lap: state.Lap{
			LapNumber: state.LapNumber((uint16(b[6]) * 256) + uint16(b[5])),
			RaceTimer: state.RaceTimer(rtd),
			LapTime:   state.LapTime(time.Duration((float64((uint16(b[2])*256)+uint16(b[3])) / 99.25) * float64(time.Second))),
		},
		TriggerValue: state.TriggerValue(0x7F & b[7]),
		Ontrack:      0x80&b[7] == 0x80,
	}
	return e
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
	o.state = c.state

	return []byte{
		o.state | o.settings.pitLane.lapTrigger | o.settings.pitLane.lapCounting,
		o.settings.maxSpeed,
		controller,
		cmd,
		parameter,
		0x00,     // unused
		0x00,     // unused
		timer[0], // Race timer
		timer[1], // Race timer
		timer[2], // Race timer
	}
}
