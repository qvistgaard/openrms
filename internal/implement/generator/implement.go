package generator

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"math/rand"
	"time"
)

type Generator struct {
	cars     uint8
	interval uint
	events   chan implement.Event
	race     *Race
}

func (g *Generator) Car(car types.Id) implement.CarImplementer {
	return NewCar(uint8(car))
}

func (g *Generator) Track() implement.TrackImplementer {
	return NewTrack()
}

func (g *Generator) Race() implement.RaceImplementer {
	return g.race
}

func (g *Generator) Init(ctx context.Context) {
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

func (g *Generator) EventChannel() <-chan implement.Event {
	return g.events
}

func (g *Generator) eventGenerator(carId uint8, interval uint) implement.Event {
	g.race.laps = uint16(0)
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Millisecond):

			var duration time.Duration
			if g.race.raceStatus == race.RaceRunning {
				duration = calculateRaceDuration(g.race.raceDuration, g.race.raceStart, time.Now())
			} else {
				duration = g.race.raceDuration
			}
			g.events <- implement.Event{
				RaceTimer: duration,
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
						Number:  g.race.laps,
						LapTime: time.Duration(rand.Intn(10000)) * time.Millisecond,
					},
				},
			}

		}
	}

}
