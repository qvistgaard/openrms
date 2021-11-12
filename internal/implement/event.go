package implement

import (
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type Event struct {
	RaceTimer time.Duration
	Car       Car
}

type Controller struct {
	BatteryWarning bool
	Link           bool
	TrackCall      bool
	ArrowUp        bool
	ArrowDown      bool
	TriggerValue   float64
}

type Car struct {
	Id         types.Id
	Reset      bool
	InPit      bool
	Deslotted  bool
	Controller Controller
	Lap        Lap
}

type Lap struct {
	Number  uint16
	LapTime time.Duration
}
