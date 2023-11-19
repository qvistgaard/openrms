package generator

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"math/rand"
	"time"
)

type Generator struct {
	cars     uint8
	interval uint
	events   chan implement.Event
}

func (g *Generator) Car(car types.Id) implement.CarImplementer {
	return NewCar(uint8(car))
}

func (g *Generator) Track() implement.TrackImplementer {
	return NewTrack()
}

func (g *Generator) Race() implement.RaceImplementer {
	return NewRace()
}

func (g *Generator) Init(ctx context.Context, processor reactive.ValuePostProcessor) {
	// NOOP
}

func (g *Generator) EventLoop() error {
	for i := 1; i <= int(g.cars); i++ {
		go g.eventGenerator(uint8(i), g.interval)
	}
	for {
		select {
		case <-time.After(10 * time.Second):
		}
	}
}

func (g *Generator) EventChannel() <-chan implement.Event {
	return g.events
}

func (g *Generator) eventGenerator(carId uint8, interval uint) implement.Event {
	laps := uint16(0)
	start := time.Now()
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Millisecond):
			laps++
			g.events <- implement.Event{
				RaceTimer: time.Now().Sub(start),
				Car: implement.Car{
					Id:        types.IdFromUint(carId),
					Reset:     false,
					InPit:     false,
					Deslotted: true,
					Controller: implement.Controller{
						BatteryWarning: false,
						Link:           false,
						TrackCall:      false,
						ArrowUp:        false,
						ArrowDown:      false,
						TriggerValue:   rand.Float64() * 100,
					},

					Lap: implement.Lap{
						Number:  laps,
						LapTime: time.Duration(rand.Intn(10000)) * time.Millisecond,
					},
				},
			}

		}
	}

}
