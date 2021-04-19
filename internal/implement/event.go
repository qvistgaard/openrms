package implement

import (
	"github.com/qvistgaard/openrms/internal/state"
	"time"
)

type Command interface {
	Exectute()
}

type CommandObject struct {
	Id      uint8
	Changes changes
}

type changes struct {
	Race map[string]state.StateInterface
	Car  map[string]state.StateInterface
}

type EventInterface interface {
}

type Event struct {
	Id           uint8
	Controller   Controller
	Car          Car
	LapTime      time.Duration
	LapNumber    uint16
	TriggerValue uint8
	Ontrack      bool
}

type Controller struct {
	BatteryWarning bool
	Link           bool
	TrackCall      bool
	ArrowUp        bool
	ArrowDown      bool
}

type Car struct {
	Reset bool
	InPit bool
}

func (e *Event) SetCarState(c *state.Car) {
	if e.Id == c.Id() {
		c.Set(state.CarEventSequence, c.Get(state.CarEventSequence).(uint)+1)
		c.Set(state.CarOnTrack, e.Ontrack)
		c.Set(state.CarControllerLink, e.Controller.Link)
		c.Set(state.CarLapNumber, e.LapNumber)
		c.Set(state.CarLapTime, e.LapTime)
		c.Set(state.CarInPit, e.Car.InPit)
		c.Set(state.CarReset, e.Car.Reset)
		c.Set(state.ControllerTriggerValue, e.TriggerValue)
		c.Set(state.ControllerBtnUp, e.Controller.ArrowUp)
		c.Set(state.ControllerBtnDown, e.Controller.ArrowDown)
		c.Set(state.ControllerBtnTrackCall, e.Controller.TrackCall)
		c.Set(state.ControllerBatteryWarning, e.Controller.BatteryWarning)
	}
}
