package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state"
	"math/rand"
	"time"
)

type Generator struct {
	cars     uint8
	interval uint
	events   chan implement.Event
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

func (g *Generator) SendCarState(c state.CarState) error {
	return nil
}

func (g *Generator) SendRaceState(r state.CourseState) error {
	return nil
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
					Id:        carId,
					Reset:     false,
					InPit:     false,
					Deslotted: true,
					Controller: implement.Controller{
						BatteryWarning: false,
						Link:           false,
						TrackCall:      false,
						ArrowUp:        false,
						ArrowDown:      false,
						TriggerValue:   rand.Float64(),
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

func (g *Generator) ResendCarState(c *state.Car) {
}
