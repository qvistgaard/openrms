package oxigen

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/goserial"
	"go.bug.st/serial/enumerator"
	"io"
	"strings"
	"sync"
	"time"
)

type Oxigen struct {
	waitGroup sync.WaitGroup

	serial     io.ReadWriteCloser
	version    string
	commands   chan Command
	bufferSize int
	mutex      sync.Mutex
	cars       map[types.CarId]Car
	track      *Track
	race       *Race
	running    bool
	start      time.Time
}

func CreateUSBConnection(device *string) (io.ReadWriteCloser, error) {
	var oxigenPort string
	if device == nil || *device == "" {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			return nil, errors.New("no serial ports found")
		}
		for _, port := range ports {
			log.WithField("port", port.Name).
				WithField("usb", port.IsUSB).
				WithField("vendor", port.VID).
				WithField("product", port.PID).
				WithField("name", port.Product).
				Debug("found COM port")
			if port.IsUSB && strings.ToUpper(port.VID) == "1FEE" && port.PID == "0002" {
				oxigenPort = port.Name
				log.WithField("port", oxigenPort).
					WithField("vendor", port.VID).
					WithField("product", port.PID).
					Info("oxigen COM port identified")
			}
		}
		if oxigenPort == "" {
			return nil, errors.New("oxigen dongle not found")
		}
	} else {
		oxigenPort = *device
	}
	log.WithField("port", oxigenPort).Info("Using oxigenport")

	options := &serial.Config{Name: oxigenPort, Baud: 115200, ReadTimeout: time.Millisecond * 20}

	// Open the port.
	port, err := serial.OpenPort(options)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func CreateImplement(serial io.ReadWriteCloser) (*Oxigen, error) {
	var err error
	o := &Oxigen{
		serial:     serial,
		commands:   make(chan Command, 1024),
		bufferSize: 1024,
		cars:       make(map[types.CarId]Car),
		track:      NewTrack(),
		race:       NewRace(),
		start:      time.Now(),
	}

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		err := o.serial.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	versionResponse := make([]byte, 5)
	_, err = io.ReadFull(o.serial, versionResponse)
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
	time.Sleep(1000 * time.Millisecond)

	return o, nil
}

func (o *Oxigen) Car(car types.CarId) drivers.Car {
	return NewCar(o, car)
}

func (o *Oxigen) Track() drivers.Track {
	return o.track
}

func (o *Oxigen) Race() drivers.Race {
	return o.race
}

func (o *Oxigen) Start(c chan<- drivers.Event) error {
	o.running = true
	o.start = time.Now()
	go o.dataExchangeLoop(c)
	return nil
}

func (o *Oxigen) Stop() error {
	if !o.running {
		return errors.New("driver not running")
	}
	o.running = false
	o.waitGroup.Wait()
	return nil
}

func (o *Oxigen) dataExchangeLoop(c chan<- drivers.Event) {
	o.waitGroup.Add(1)

	defer o.cleanup()
	defer o.waitGroup.Done()

	for o.running {
		bytesReceived := 0

		for bytesReceived == 0 {
			err := o.tx()
			if err != nil {
				log.Error(err)
			}

			bytesReceived, err = o.rx(c)
			if err != nil {
				bytesReceived = 0
				// log.Trace(err)
			}
			time.Sleep(250 * time.Millisecond)

		}
	}
}

func (o *Oxigen) cleanup() {
	err := o.serial.Close()
	if err != nil {
		return
	}
	log.Error(err)
}

func (o *Oxigen) rx(c chan<- drivers.Event) (int, error) {
	buffer := make([]byte, 13)

	r := io.LimitReader(o.serial, 13)
	var err error

	n, err := r.Read(buffer)
	if err == nil {
		if n == 13 {
			log.WithField("message", fmt.Sprintf("%x", buffer)).
				WithField("bytes", n).
				Trace("received message from oxygen dongle")
			o.event(c, buffer)
			return n, nil
		} else if n > 0 {
			log.WithField("message", fmt.Sprintf("%x", buffer)).
				WithField("bytes", n).
				Trace("message with incorrect byte count received from oxygen dongle")
			return n, errors.New("message with incorrect byte count received from oxygen dongle")
		}
	}
	return n, err
}

func (o *Oxigen) sendCarCommand(car uint8, code byte, value uint8) {
	command := newCommand(&car, code, value)
	o.commands <- command
}

func (o *Oxigen) tx() error {
	select {
	case cmd := <-o.commands:
		if float32(len(o.commands)) > (float32(o.bufferSize) * 0.9) {
			log.WithFields(map[string]interface{}{
				"bufferSize": o.bufferSize,
				"size":       len(o.commands),
			}).Warn("too many commands on code buffer")
		} else {
			log.WithFields(map[string]interface{}{
				"bufferSize": o.bufferSize,
				"size":       len(o.commands),
			}).Trace("Buffer size")

		}

		b := o.packCommand(cmd, packRaceCounter(o.start))
		if cmd.id != nil {
			log.WithFields(map[string]interface{}{
				"message": fmt.Sprintf("%x", b),
				"car":     fmt.Sprintf("%x", *cmd.id),
				"value":   fmt.Sprintf("%x", cmd.value),
				"cmd":     fmt.Sprintf("%x", cmd.code),
				"hex":     fmt.Sprintf("%s", hex.Dump(b)),
			}).Debug("sending message to oxygen dongle")
		}

		l, err := o.serial.Write(b)
		if err != nil {
			panic(err)
		}
		log.WithFields(map[string]interface{}{
			"message": fmt.Sprintf("%x", b),
			"size":    fmt.Sprintf("%d", l),
			"buffer":  len(o.commands),
		}).Trace("send message to oxygen dongle")
		if err != nil {
			return err
		}
	case <-time.After(100 * time.Millisecond):
		if len(o.commands) == 0 {
			log.Trace("oxigen: sending keep-alive")
			o.commands <- newEmptyCommand()
		}
		// time.Sleep(400 * time.Millisecond)
	case <-time.After(5000 * time.Millisecond):
		if !o.running {
			return nil
		}
	}
	return nil
}

func (o *Oxigen) event(c chan<- drivers.Event, b []byte) {
	rt := unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4])
	lt := unpackLapTime(b[2], b[3])
	lapNumber := (uint32(b[6]) * 256) + uint32(b[5])

	car := NewCar(o, types.IdFromUint(b[1]))
	c <- events.NewControllerTriggerValueEvent(car, float64(0x7F&b[7]))
	c <- events.NewControllerTrackCallButton(car, 0x08&b[0] == 0x08)
	c <- events.NewLap(car, lapNumber, lt, rt)
	c <- events.NewInPit(car, 0x40&b[8] == 0x40)
	// TODO: Deprecate and remove
	c <- events.NewDeslotted(car, !(0x80&b[7] == 0x80))
	c <- events.NewOnTrack(car, 0x80&b[7] == 0x80)

	// e := drivers.GenericEvent(nil)
	/*
		NOTE: Keep this here for future reference

					drivers.Event{
				RaceTimer: unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4]),
				Car: drivers.Car{
					CarId:        types.IdFromUint(b[1]), // OK
					Reset:     0x01&b[0] == 0x01, //
					InPit:     0x40&b[8] == 0x40, // OK
					Deslotted: !(0x80&b[7] == 0x80), // OK
					Controller: drivers.Controller{
						BatteryWarning: 0x04&b[0] == 0x04,
						Link:           0x02&b[0] == 0x02,
						TrackCall:      0x08&b[0] == 0x08, // OK
						ArrowUp:        0x20&b[0] == 0x20,
						ArrowDown:      0x40&b[0] == 0x40,
						TriggerValue:   float64(0x7F & b[7]), // OK
					},
					Lap: drivers.Lap{
						Number:  (uint16(b[6]) * 256) + uint16(b[5]), // OK
						LapTime: unpackLapTime(b[2], b[3]), // OK
					},
				},
			}
	*/
	// return e
}

func (o *Oxigen) packCommand(c Command, timer []byte) []byte {
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
	be := binary.BigEndian.Uint64([]byte{0, 0, 0, 0, 0, 0, high, low})
	lt := float64(be) / 99.25
	ltd := time.Duration(lt * float64(time.Second))
	return ltd
}

func unpackRaceTime(b [4]byte, lag byte) time.Duration {
	be := make([]byte, 8)
	copy(be[4:], b[:])

	rt := binary.BigEndian.Uint64(be) - uint64(lag)
	// rt := (uint(b[0]) * 16777216) + (uint(b[1]) * 65536) + (uint(b[2]) * 256) + uint(b[3]) - uint(lag)
	return time.Duration(rt*10) * time.Millisecond
}

func packRaceCounter(start time.Time) []byte {
	centiSeconds := time.Now().Sub(start).Milliseconds() / 10
	be := make([]byte, 8)
	binary.BigEndian.PutUint64(be, uint64(centiSeconds))
	return be[len(be)-3:]
}

func unpackRaceCounter(b [3]byte) time.Duration {
	be := make([]byte, 8)
	copy(be[5:], b[:])

	u := binary.BigEndian.Uint64(be)
	return time.Duration(u*10) * time.Millisecond
}
