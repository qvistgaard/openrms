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

func (g *Generator) SendCarState(c state.CarChanges) error {
	return nil
}

func (g *Generator) SendRaceState(r state.CourseChanges) error {
	return nil
}

func (g *Generator) eventGenerator(car uint8, interval uint) implement.Event {
	laps := state.LapNumber(0)
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Millisecond):
			laps++
			g.events <- implement.Event{
				Id: state.CarId(car),
				Controller: implement.Controller{
					BatteryWarning: false,
					Link:           false,
					TrackCall:      false,
					ArrowUp:        false,
					ArrowDown:      false,
				},
				Car: implement.Car{
					Reset: false,
					InPit: false,
				},
				Lap: state.Lap{
					LapNumber: laps,
					RaceTimer: 0,
					LapTime:   state.LapTime(time.Duration(rand.Intn(10000)) * time.Millisecond),
				},
				TriggerValue: state.TriggerValue(uint8(rand.Int31())),
				Ontrack:      true,
			}

		}
	}

}