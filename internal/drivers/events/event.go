package events

import (
	"github.com/qvistgaard/openrms/internal/drivers"
)

type Event interface {
	Car() drivers.Car
}

type GenericEvent struct {
	car drivers.Car
}

func (g GenericEvent) Car() drivers.Car {
	return g.car
}

type ControllerBatteryWarning interface {
	drivers.Event
	BatteryWarning() bool
}

type ControllerLinkEvent interface {
	drivers.Event
	Link() bool
}

/*
type Event struct {
	RaceTimer time.Duration
	Car       Car
}
*/
/*
	type Controller struct {
		BatteryWarning bool
		Link           bool
		TrackCall      bool
		ArrowUp        bool
		ArrowDown      bool
		TriggerValue   float64
	}
*/
/*
type Car struct {
	Id         types.Id
	Reset      bool
	InPit      bool
	Deslotted  bool
	Controller Controller
	Lap        Lap
}
*/
/*
type Lap struct {
	Number  uint16
	LapTime time.Duration
}
*/
