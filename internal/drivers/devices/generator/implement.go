package generator

import (
	"context"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type Generator struct {
	cars     uint8
	interval uint
	events   chan drivers.Event
	race     *Race
}

func (g *Generator) Car(car types.Id) drivers.Car {
	return NewCar(uint8(car.ToUint()), 0)
}

func (g *Generator) Track() drivers.Track {
	return NewTrack()
}

func (g *Generator) Race() drivers.Race {
	return g.race
}

func (g *Generator) Init(_ context.Context) {
	// NOOP
}

func (g *Generator) EventLoop() error {
	for i := 1; i <= int(g.cars); i++ {
		go g.eventGenerator(uint8(i), g.interval)
	}
	for {
		select {
		case <-time.After(time.Duration(g.interval) * time.Millisecond):
			g.race.laps++
		}
	}
}

func (g *Generator) EventChannel() <-chan drivers.Event {
	return g.events
}

func (g *Generator) eventGenerator(carId uint8, interval uint) drivers.Event {
	g.race.laps = uint16(0)
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Millisecond):
			car := NewCar(carId, g.race.laps)

			g.events <- drivers.GenericEvent(car)
			/*
					drivers.Event{
				RaceTimer: time.Now().Sub(g.race.raceStart),
				Car: drivers.Car{
					Id:        types.IdFromUint(carId),
					Reset:     false,
					InPit:     false,
					Deslotted: false,
					Controller: drivers.Controller{
						BatteryWarning: false,
						Link:           false,
						TrackCall:      false,
						ArrowUp:        false,
						ArrowDown:      false,
						TriggerValue:   rand.Float64() * 100,
					},

					Lap: drivers.Lap{
						Number:  g.race.laps,
						LapTime: time.Duration(rand.Intn(10000)) * time.Millisecond,
					},
				},

			*/

		}
	}
}
