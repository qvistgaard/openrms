package v3

import (
	"errors"
	"fmt"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/rs/zerolog"
	"go.bug.st/serial"
	"math/rand"
	"syscall"
	"time"
)

type Driver3x struct {
	serial       serial.Port
	start        time.Time
	links        map[types.CarId]controllerLink
	cars         map[types.CarId]*Car
	expire       chan types.CarId
	readInterval int
	version      string
	tx           chan Command
	race         *Race
	track        *Track
	logger       zerolog.Logger
}

var noLinksError = errors.New("empty message from dongle and no links available")

type dongleRxMessage [13]byte

func CreateDriver(logger zerolog.Logger, connection serial.Port) (drivers.Driver, error) {
	connection.SetReadTimeout(50 * time.Millisecond)
	o := &Driver3x{
		logger: logger,
		serial: connection,
		track:  NewTrack(logger),
		race:   NewRace(logger),
		start:  time.Now(),
		links:  make(map[types.CarId]controllerLink),
		cars:   make(map[types.CarId]*Car),
		expire: make(chan types.CarId),
		tx:     make(chan Command, 1024),

		readInterval: 320,
	}
	return o, nil
}

func (d *Driver3x) Start(events chan<- drivers.Event) error {
	d.start = time.Now()
	go d.communicationLoop(events)
	go d.linkUpdateLoop(events)
	return nil
}

func (d *Driver3x) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver3x) Car(car types.CarId) drivers.Car {
	if _, ok := d.cars[car]; !ok {
		d.cars[car] = newCar(d, car)
	}
	return d.cars[car]
}

func (d *Driver3x) Track() drivers.Track {
	return d.track
}

func (d *Driver3x) Race() drivers.Race {
	return d.race
}

func (d *Driver3x) linkUpdateLoop(e chan<- drivers.Event) {
	idleDuration := 100 * time.Millisecond
	idleDelay := time.NewTimer(idleDuration)

	for {
		idleDelay.Reset(idleDuration)

		select {
		case link := <-d.expire:
			d.removeLink(link)
			e <- events.NewEnabled(newCar(d, link), false)
		case <-idleDelay.C:
			if len(d.tx) == 0 {
				d.sendStoredCarState()
			}
		}
	}

}

func (d *Driver3x) communicationLoop(events chan<- drivers.Event) {
	idleDuration := 50 * time.Millisecond
	idleDelay := time.NewTimer(idleDuration)

	for {
		idleDelay.Reset(idleDuration)
		select {
		case command := <-d.tx:
			d.writeAndRead(command, events)
		case <-idleDelay.C:
			if len(d.tx) == 0 {
				d.tx <- newEmptyCommand()
			}
		}
	}
}

func (d *Driver3x) removeLink(link types.CarId) {
	delete(d.links, link)
}

func (d *Driver3x) writeAndRead(command Command, events chan<- drivers.Event) {
	for {
		_, err := d.write(command)
		if err != nil {
			d.logger.Err(err).Msg("Failed to write command to dongle")
			continue
		} else {
			break
		}
	}
	for {
		if err := d.serial.Drain(); err != nil {
			var errno syscall.Errno
			if errors.As(err, &errno) && errors.Is(errno, syscall.EINTR) {
				d.logger.Warn().Err(err).Msg("Failed to drain after writing command. retrying...")
				continue
			}
			d.logger.Warn().Err(err).Msg("Failed to drain after writing command")
		} else {
			break
		}
	}

	for {
		read, err := d.Read()
		if errors.Is(err, errors.New("EOF")) {
			continue
		} else if errors.Is(noLinksError, err) {
			return
		} else if err != nil {
			return
		}

		if err == nil || len(read) > 0 {
			for _, slice := range read {
				d.updateLink(events, slice)
				d.event(events, slice)
			}
			return
		}
	}
}

func (d *Driver3x) write(command Command) (int, error) {
	timer := packRaceCounter(d.start)
	pack := command.pack(timer, d.race, d.track)

	n, err := d.serial.Write(pack)
	return n, err
}

func (d *Driver3x) Read() ([]dongleRxMessage, error) {
	var messages []byte
	buffer := make([]byte, 52)

	for len(messages) == 0 || len(messages)%13 != 0 {
		n, err := d.serial.Read(buffer)
		if err != nil {
			d.logger.Err(err).Msg("Failed to read buffer")
			return nil, err
		}
		if n == 0 {
			if len(d.links) == 0 {
				return []dongleRxMessage{}, noLinksError
			}
			return []dongleRxMessage{}, errors.New("empty message from dongle")
		}

		messages = append(messages, buffer[:n]...)
		// buffer = nil
	}

	if d.logger.Trace().Enabled() {
		d.logger.Trace().
			Str("messages", fmt.Sprintf("%v", messages)).
			Int("bytes", len(messages)).
			Msg("received message from dongle")
	}
	return d.splitMessages(messages), nil
}

func (d *Driver3x) event(c chan<- drivers.Event, b dongleRxMessage) {
	id := types.IdFromUint(b[1])

	rt := unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4])
	lt := unpackLapTime(b[2], b[3])
	lapNumber := (uint32(b[6]) * 256) + uint32(b[5])

	car := d.Car(id)

	c <- events.NewInPit(car, unpackPitStatus(b))
	c <- events.NewControllerTriggerValueEvent(car, float64(0x7F&b[7]))
	c <- events.NewControllerTrackCallButton(car, 0x08&b[0] == 0x08)
	c <- events.NewLap(car, lapNumber, lt, rt)
	// TODO: Deprecate and remove
	c <- events.NewDeslotted(car, !(0x80&b[7] == 0x80))
	c <- events.NewOnTrack(car, 0x80&b[7] == 0x80)
}

func (d *Driver3x) splitMessages(messages []byte) []dongleRxMessage {
	count := len(messages) / 13
	messageSlice := messages[:]
	splitSlices := make([]dongleRxMessage, count)

	for i := 0; i < count; i++ {
		copy(splitSlices[i][:], messageSlice[i*13:(i+1)*13])
	}
	return splitSlices[:count]
}

func (d *Driver3x) updateLink(e chan<- drivers.Event, message [13]byte) {
	linkId := types.IdFromUint(message[1])
	if linkId > 0 {
		if _, ok := d.links[linkId]; !ok {
			d.links[linkId] = controllerLink{
				logger: d.logger,
				id:     linkId,
				expire: d.expire,
				renew:  make(chan bool),
			}
			d.readInterval = 320
			l := d.links[linkId]
			go l.timeout()
			e <- events.NewEnabled(newCar(d, linkId), true)
		}
		d.links[linkId].renew <- true
	}
}

func (d *Driver3x) sendCarCommand(car uint8, code byte, value uint8) {
	command := newCommand(&car, code, value)
	d.tx <- command
}

func (d *Driver3x) sendStoredCarState() {
	if len(d.cars) == 0 {
		return // Early return if there are no cars
	}

	var randomCarId types.CarId
	var i int
	randomInt := rand.Intn(len(d.cars))
	for id := range d.cars {
		if i == randomInt {
			randomCarId = id
			break
		}
		i++
	}

	car := d.cars[randomCarId]
	switch rand.Intn(4) {
	case 0:
		car.sendMaxBreaking()
	case 1:
		car.sendMinSpeed()
	case 2:
		car.sendMaxSpeed()
	case 3:
		car.sendPitLaneMaxSpeed()
	}

}
