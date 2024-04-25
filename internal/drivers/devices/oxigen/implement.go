package oxigen

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	v3 "github.com/qvistgaard/openrms/internal/drivers/devices/oxigen/v3"
	"github.com/rs/zerolog"
	"go.bug.st/serial"
	"io"
	"time"
)

func CreateImplement(logger zerolog.Logger, connection serial.Port) (drivers.Driver, error) {
	var err error
	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = connection.Write(versionRequest)
	if err != nil {
		err := connection.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	versionResponse := make([]byte, 5)
	_, err = io.ReadFull(connection, versionResponse)
	v, _ := version.NewVersion(fmt.Sprintf("%d.%d", versionResponse[0], versionResponse[1]))
	constraint, _ := version.NewConstraint(">= 3.10")

	logger.Trace().
		Str("message", fmt.Sprintf("%x", versionResponse)).
		Int("bytes", len(versionResponse)).
		Msg("received message from oxygen dongle")

	if !constraint.Check(v) {
		return nil, errors.New(fmt.Sprintf("Unsupported dongle version: %s", v))
	}

	logger.Info().
		Stringer("version", v).
		Msgf("Connected to oxigen dongle. Dongle version: %s", v)
	time.Sleep(1000 * time.Millisecond)

	return v3.CreateDriver(logger, connection)
}

/*
type Oxigen struct {
	waitGroup sync.WaitGroup

	serial       io.ReadWriteCloser
	version      string
	commands     chan oxigen.Command
	bufferSize   int
	mutex        sync.Mutex
	cars         map[types.CarId]Car
	track        *Track
	race         *Race
	running      bool
	start        time.Time
	links        map[types.CarId]v3.controllerLink
	expire       chan types.CarId
	readInterval int
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

	options := &serial.Config{Name: oxigenPort, Baud: 115200, ReadTimeout: time.Millisecond * 50}

	// Open the port.
	port, err := serial.OpenPort(options)
	if err != nil {
		return nil, err
	}
	return port, nil
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

	defer log.Error("DEL died")
	defer o.cleanup()
	defer o.waitGroup.Done()

	for o.running {
		bytesReceived := 0
		nextMessage, err := o.tx()
		if err != nil {
			panic(err)
		}

		for bytesReceived == 0 {
			err = o.sendMessage(nextMessage)
			if err != nil {
				log.Error(err)
			} else {
				o.rxCoolDown()
				bytesReceived, err = o.rx(c)
				if err != nil || bytesReceived == 0 {
					if len(o.links) > 0 {
						o.readInterval = o.readInterval + 10
						log.WithField("interval", o.readInterval).Warn("Read error from dongle. adjusting cooldown interval")
					} else {
						break
					}
					bytesReceived = 0
				}
			}
		}
	}
}

func (o *Oxigen) sendMessage(nextMessage []byte) error {
	l, err := o.serial.Write(nextMessage)
	if err != nil {
		return err
	}
	log.WithFields(map[string]interface{}{
		"message": fmt.Sprintf("%x", nextMessage),
		"size":    fmt.Sprintf("%d", l),
		"buffer":  len(o.commands),
	}).Trace("send message to oxygen dongle")
	return nil
}

func (o *Oxigen) rxCoolDown() {
	time.Sleep(time.Duration(o.readInterval/(len(o.links)+1)) * time.Millisecond)
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
	command := oxigen.newCommand(&car, code, value)
	o.commands <- command
}

func (o *Oxigen) tx() ([]byte, error) {
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

		return b, nil

	case <-time.After(10 * time.Millisecond):
		if len(o.commands) == 0 {
			log.Trace("oxigen: sending keep-alive")
			o.commands <- oxigen.newEmptyCommand()
		}
		// time.Sleep(400 * time.Millisecond)
	case <-time.After(5000 * time.Millisecond):
		if !o.running {
			return nil, nil
		}
	}
	return nil, nil
}

func (o *Oxigen) event(c chan<- drivers.Event, b []byte) {
	id := types.IdFromUint(b[1])

	rt := unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4])
	lt := unpackLapTime(b[2], b[3])
	lapNumber := (uint32(b[6]) * 256) + uint32(b[5])

	car := NewCar(o, id)

	c <- events.NewControllerTriggerValueEvent(car, float64(0x7F&b[7]))
	c <- events.NewControllerTrackCallButton(car, 0x08&b[0] == 0x08)
	c <- events.NewLap(car, lapNumber, lt, rt)
	c <- events.NewInPit(car, 0x40&b[8] == 0x40)
	// TODO: Deprecate and remove
	c <- events.NewDeslotted(car, !(0x80&b[7] == 0x80))
	c <- events.NewOnTrack(car, 0x80&b[7] == 0x80)
*/
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
/*
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
*/
