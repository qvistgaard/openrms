package generator

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"math/rand"
	"sync"
	"time"
)

type Generator struct {
	waitGroup sync.WaitGroup
	started   bool
	cars      uint8
	interval  uint
	race      *Race
}

func (g *Generator) Start(c chan<- drivers.Event) error {
	g.started = true

	for i := 1; i <= int(g.cars); i++ {
		go g.eventGenerator(types.CarId(i), c, g.interval)
		time.Sleep(250 * time.Millisecond)
	}
	go g.updateLapCounter()
	return nil
}

func (g *Generator) Stop() error {
	if !g.started {
		return errors.New("driver not running")
	}
	g.started = false
	g.waitGroup.Wait()
	return nil
}

func (g *Generator) Car(car types.CarId) drivers.Car {
	return NewCar(car, 0)
}

func (g *Generator) Track() drivers.Track {
	return NewTrack()
}

func (g *Generator) Race() drivers.Race {
	return g.race
}

func (g *Generator) eventGenerator(carId types.CarId, c chan<- drivers.Event, interval uint) {
	g.waitGroup.Add(1)
	g.race.laps = uint32(0)

	defer g.waitGroup.Done()
	for g.started {
		select {
		case <-time.After(time.Duration(interval) * time.Millisecond):
			car := NewCar(carId, g.race.laps)
			deslot := false // rand.Float32() < 0.07
			if deslot {
				c <- events.NewOnTrack(car, false)
				c <- events.NewDeslotted(car, true)
			} else {
				c <- events.NewOnTrack(car, true)
				c <- events.NewDeslotted(car, false)
			}
			if g.race.raceStatus == race.Running {
				c <- events.NewControllerTriggerValueEvent(car, float64(100))
				c <- events.NewLap(car, g.race.laps, time.Duration(rand.Intn(10000))*time.Millisecond, time.Now().Sub(g.race.raceStart))
			}
		}
	}
}

func (g *Generator) updateLapCounter() {
	g.started = true
	g.waitGroup.Add(1)

	defer g.waitGroup.Done()
	for g.started {
		select {
		case <-time.After(time.Duration(g.interval) * time.Millisecond):
			if g.race.raceStatus == race.Running {
				g.race.laps++
			}
		}
	}
}
