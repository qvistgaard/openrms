package v3

import (
	"errors"
	"fmt"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

type Driver3x struct {
	serial       io.ReadWriteCloser
	start        time.Time
	links        map[types.CarId]controllerLink
	expire       chan types.CarId
	readInterval int
	version      string
	tx           chan Command
	race         *Race
	track        *Track
}

type dongleRxMessage [13]byte

func CreateDriver(serial io.ReadWriteCloser) (drivers.Driver, error) {
	o := &Driver3x{
		serial: serial,
		track:  NewTrack(),
		race:   NewRace(),
		start:  time.Now(),
		links:  make(map[types.CarId]controllerLink),
		expire: make(chan types.CarId),
		tx:     make(chan Command, 1024),

		readInterval: 100,
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
	return newCar(d, car)

}

func (d *Driver3x) Track() drivers.Track {
	return d.track
}

func (d *Driver3x) Race() drivers.Race {
	return d.race
}

func (d *Driver3x) linkUpdateLoop(e chan<- drivers.Event) {
	for {
		select {
		case link := <-d.expire:
			d.removeLink(link)
			e <- events.NewEnabled(newCar(d, link), false)
		}
	}

}

func (d *Driver3x) communicationLoop(events chan<- drivers.Event) {
	for {
		select {
		case command := <-d.tx:
			d.writeAndRead(command, events)
		case <-time.After(50 * time.Millisecond):
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
			log.Error("Failed to write command to dongle", err)
			return
		}

		read, err := d.Read()
		if err != nil && read == nil {
			if err.Error() == "EOF" {
				return
			}
			log.Error("Failed to read from buffer", err)
			return
		}
		if err != nil || len(read) > 0 {
			for _, slice := range read {
				/*				fmt.Printf("%v ", 0x40&slice[8] == 0x40)
								fmt.Printf("%08b %08b ", slice[4], slice[8])
								fmt.Printf("Slice %d: %v\n", i, slice)*/
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

	field := log.WithField("message", fmt.Sprintf("%v", pack)).
		WithField("bytes", n).
		WithField("car", command.id)

	if command.id != nil {
		field.Debug("send message to dongle")
	} else {
		field.Trace("send message to dongle")
	}

	return n, err
}

func (d *Driver3x) Read() ([]dongleRxMessage, error) {
	var messages []byte
	buffer := make([]byte, 52)

	for len(messages) == 0 || len(messages)%13 != 0 {
		time.Sleep(time.Duration(d.readInterval) / time.Duration(len(d.links)+1) * time.Millisecond)

		n, err := d.serial.Read(buffer)
		if err != nil {
			return nil, err
		}

		log.WithField("message", fmt.Sprintf("%v", buffer)).
			WithField("bytes", n).
			Trace("received message from dongle")

		if n == 0 {
			if len(d.links) == 0 {
				return []dongleRxMessage{}, errors.New("empty message from dongle and no links available")
			}
			d.readInterval = d.readInterval + 10
			log.WithField("interval", d.readInterval).Error("Read timeout, increasing read interval")
			return []dongleRxMessage{}, errors.New("empty message from dongle")
		}

		messages = append(messages, buffer[:n]...)
	}
	return d.splitMessages(messages), nil
}

func (d *Driver3x) event(c chan<- drivers.Event, b dongleRxMessage) {
	id := types.IdFromUint(b[1])

	rt := unpackRaceTime([4]byte{b[9], b[10], b[11], b[12]}, b[4])
	lt := unpackLapTime(b[2], b[3])
	lapNumber := (uint32(b[6]) * 256) + uint32(b[5])

	car := newCar(d, id)

	c <- events.NewControllerTriggerValueEvent(car, float64(0x7F&b[7]))
	c <- events.NewControllerTrackCallButton(car, 0x08&b[0] == 0x08)
	c <- events.NewLap(car, lapNumber, lt, rt)
	c <- events.NewInPit(car, 0x40&b[8] == 0x40)
	// TODO: Deprecate and remove
	c <- events.NewDeslotted(car, !(0x80&b[7] == 0x80))
	c <- events.NewOnTrack(car, 0x80&b[7] == 0x80)
}

func (d *Driver3x) splitMessages(messages []byte) []dongleRxMessage {
	count := len(messages) / 13
	messageSlice := messages[:]
	var splitSlices [4]dongleRxMessage

	for i := 0; i < count; i++ {
		copy(splitSlices[i][:], messageSlice[i*13:(i+1)*13])
	}
	return splitSlices[:count]
}

func (d *Driver3x) updateLink(e chan<- drivers.Event, message [13]byte) {
	linkId := types.IdFromUint(message[1])

	if _, ok := d.links[linkId]; !ok {
		d.links[linkId] = controllerLink{
			id:     linkId,
			expire: d.expire,
			renew:  make(chan bool),
		}
		d.readInterval = 100
		l := d.links[linkId]
		go l.timeout()
		e <- events.NewEnabled(newCar(d, linkId), true)
	}
	d.links[linkId].renew <- true
}

func (d *Driver3x) sendCarCommand(car uint8, code byte, value uint8) {
	command := newCommand(&car, code, value)
	d.tx <- command
}
