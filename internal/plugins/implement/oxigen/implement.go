package oxigen

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"io"
	"sync"
	"time"
)

type Oxigen struct {
	serial     io.ReadWriteCloser
	version    string
	commands   chan Command
	events     chan implement.Event
	bufferSize int
	timer      []byte
	buffer     *bufio.Reader
	mutex      sync.Mutex
	cars       map[types.Id]Car
	track      *Track
	race       *Race
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
	o := &Oxigen{
		serial:     serial,
		commands:   make(chan Command, 1024),
		events:     make(chan implement.Event, 1024),
		buffer:     bufio.NewReaderSize(serial, 1024),
		bufferSize: 1024,
		cars:       make(map[types.Id]Car),
		track:      NewTrack(),
		race:       NewRace(),
	}

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		o.serial.Close()
		return nil, err
	}

	versionResponse := make([]byte, 5)
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

func (o *Oxigen) Car(car types.Id) implement.CarImplementer {
	return NewCar(o, uint8(car))
}

func (o *Oxigen) Track() implement.TrackImplementer {
	return o.track
}

func (o *Oxigen) Race() implement.RaceImplementer {
	return o.race
}

func (o *Oxigen) Init(ctx context.Context, processor reactive.ValuePostProcessor) {
	//	o.track.Init(ctx, processor)
	//	o.race.Init(ctx, processor)
}

func (o *Oxigen) EventLoop() error {
	defer func() {
		o.serial.Close()
		panic("oxigen: keep-alive routine failed.")
	}()

	o.timer = []byte{0x00, 0x00, 0x00}

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

func (o *Oxigen) sendCarCommand(car *uint8, code byte, value uint8) {
	command := newCommand(car, code, value)
	o.commands <- command
}

func (o *Oxigen) sendCommand() {
	defer func() {
		panic("oxigen: keep-alive routine failed.")
	}()
	for {
		o.mutex.Lock()
		select {
		case cmd := <-o.commands:
			if float32(len(o.commands)) > (float32(o.bufferSize) * 0.9) {
				log.WithFields(map[string]interface{}{
					"bufferSize": o.bufferSize,
					"size":       len(o.commands),
				}).Warn("too many commands on code buffer")
			}

			b := o.command(cmd, o.timer)
			if cmd.id != nil {
				log.WithFields(map[string]interface{}{
					"message": fmt.Sprintf("%x", b),
					"car":     fmt.Sprintf("%x", *cmd.id),
					"value":   fmt.Sprintf("%x", cmd.value),
					"cmd":     fmt.Sprintf("%x", cmd.code),
					"hex":     fmt.Sprintf("%s", hex.Dump(b)),
				}).Trace("sendind message to oxygen dongle")
			}
			l, err := o.serial.Write(b)
			if err != nil {
				panic(err)
			}
			log.WithFields(map[string]interface{}{
				"message": fmt.Sprintf("%x", b),
				"size":    fmt.Sprintf("%d", l),
			}).Trace("send message to oxygen dongle")
			time.Sleep(10 * time.Millisecond)
		}
		o.mutex.Unlock()
	}
}

func (o *Oxigen) keepAlive() {
	defer func() {
		panic("oxigen: keep-alive routine failed.")
	}()
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			if len(o.commands) == 0 {
				o.commands <- newEmptyCommand()
				log.Trace("oxigen: sent keep-alive")
			}
		}
	}
}

func (o *Oxigen) EventChannel() <-chan implement.Event {
	return o.events
}

func (o *Oxigen) event(b []byte) implement.Event {
	e := implement.Event{
		RaceTimer: unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4]),
		Car: implement.Car{
			Id:        types.IdFromUint(b[1]),
			Reset:     0x01&b[0] == 0x01,
			InPit:     0x40&b[8] == 0x40,
			Deslotted: !(0x80&b[7] == 0x80),
			Controller: implement.Controller{
				BatteryWarning: 0x04&b[0] == 0x04,
				Link:           0x02&b[0] == 0x02,
				TrackCall:      0x08&b[0] == 0x08,
				ArrowUp:        0x20&b[0] == 0x20,
				ArrowDown:      0x40&b[0] == 0x40,
				TriggerValue:   float64(0x7F & b[7]),
			},
			Lap: implement.Lap{
				Number:  (uint16(b[6]) * 256) + uint16(b[5]),
				LapTime: unpackLapTime(b[2], b[3]),
			},
		},
	}
	return e
}

func (o *Oxigen) command(c Command, timer []byte) []byte {
	var cmd byte = 0x00
	var parameter byte = 0x00
	var controller byte = 0x00

	cmd = c.code
	parameter = c.value
	if c.id != nil {
		controller = *c.id
		cmd = 0x80 | cmd
	} else {
		controller = 0x00
		cmd = 0x00 | cmd
	}

	return []byte{
		o.race.status | o.track.pitLane.lapCounting | o.track.pitLane.lapCountingOption,
		o.track.maxSpeed,
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

func unpackLapTime(high byte, low byte) time.Duration {
	lt := float64((uint(high)*256)+uint(low)) / 99.25
	ltd := time.Duration(lt * float64(time.Second))
	return ltd
}

func unpackRaceTime(b [4]byte, lag byte) time.Duration {
	rt := (uint(b[0]) * 16777216) + (uint(b[1]) * 65536) + (uint(b[2]) * 256) + uint(b[3]) - uint(lag)
	return time.Duration(rt*10) * time.Millisecond
}
