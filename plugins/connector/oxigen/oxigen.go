package oxigen

import (
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

func Connect(serial ReadWriteCloserConnector) (*Oxigen, error) {
	var err error
	o := new(Oxigen)

	o.serial, err = serial.connect()
	if err != nil {
		return nil, err
	}

	versionRequest := []byte{0x06, 0x06, 0x06, 0x06, 0x00, 0x00, 0x00} // Get dongle version bytecode
	_, err = o.serial.Write(versionRequest)
	if err != nil {
		o.serial.Close()
		return nil, err
	}

	versionResponse := make([]byte, 13)
	_, err = o.serial.Read(versionResponse)
	v, _ := version.NewVersion(fmt.Sprintf("%d.%d", versionResponse[0], versionResponse[1]))
	constraint, _ := version.NewConstraint(">= 3.10")

	if !constraint.Check(v) {
		return nil, errors.New(fmt.Sprintf("Unsupported dongle version: %s", v))
	}
	o.version = v.String()
	log.Printf("Connected to oxigen dongle. Dongle version: %s", v)

	return o, nil
}

func (oxigen Oxigen) Closer() error {
	return oxigen.serial.Close()
}

type Oxigen struct {
	state    byte
	serial   io.ReadWriteCloser
	settings Settings
	version  string
}

type Settings struct {
	maxSpeed byte
	pitLane  PitLane
}

type PitLane struct {
	lapCounting byte
	lapTrigger  byte
}

func (oxigen Oxigen) EventLoop(input queue.Queue, output queue.Queue) error {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	var err error
	for {
		_, _ = input.DequeueOrWaitForNextElementContext(ctx)
		b := []byte{
			oxigen.state | oxigen.settings.pitLane.lapTrigger | oxigen.settings.pitLane.lapCounting,
			oxigen.settings.maxSpeed,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
		}

		_, err = oxigen.serial.Write(b)
		if err != nil {
			log.Fatalf("port.Write: %v", err)
			break
		}
		for {
			buffer := make([]byte, 13)
			_, err = oxigen.serial.Read(buffer)
			if err != nil {
				break
			}

			err = output.Enqueue(buffer)
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func (oxigen Oxigen) maxSpeed(speed uint8) bool {
	oxigen.settings.maxSpeed = speed
	return true
}

func (oxigen Oxigen) start() bool {
	oxigen.state = 0x03
	return true
}

func (oxigen Oxigen) pitLaneLapCount(enabled bool, entry bool) bool {
	if enabled {
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

func (oxigen Oxigen) stop() bool {
	oxigen.state = 0x01
	return true
}

func (oxigen Oxigen) pause() bool {
	oxigen.state = 0x04
	return true
}

func (oxigen Oxigen) flag(lc bool) bool {
	if lc {
		oxigen.state = 0x05
	} else {
		oxigen.state = 0x15
	}
	return true
}
