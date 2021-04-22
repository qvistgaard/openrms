package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/state"
)

const (
	CarLastLaps = "car-last-laps"
)

type Board interface {
	updateCar(car *state.Car)
}

type LastLaps interface {
	update(car state.Lap) LastLapDefault
}
type LastLapDefault struct {
	Laps []state.Lap
}

type BoardDefault struct {
}

func (l *LastLapDefault) update(lap state.Lap) *LastLapDefault {

	slice := []state.Lap{lap}
	//	l.Laps =

	return &LastLapDefault{
		Laps: append(slice, l.Laps[0:len(l.Laps)-1]...),
	}
}

func (l *LastLapDefault) Compare(v state.ComparableChange) bool {
	if c, ok := v.(*LastLapDefault); ok {
		return c.Laps[0].LapNumber != l.Laps[0].LapNumber
	}
	return false
}

func (b *BoardDefault) InitializeCarState(car *state.Car) {
	car.Set(CarLastLaps, &LastLapDefault{
		Laps: make([]state.Lap, 5),
	})
	car.Subscribe(state.CarLap, b)
}

func (b *BoardDefault) InitializeRaceState(race *state.Course) {
}

func (b *BoardDefault) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if l, ok := v.Get().(state.Lap); ok && v.Name() == state.CarLap {
			last := c.Get(CarLastLaps).(*LastLapDefault)
			c.Set(CarLastLaps, last.update(l))
		}

	}
}
