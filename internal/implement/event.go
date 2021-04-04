package implement

import (
	"github.com/qvistgaard/openrms/internal/state"
	"time"
)

type Command struct {
	Id      uint8
	Changes changes
}

type changes struct {
	Race map[string]state.StateInterface
	Car  map[string]state.StateInterface
}

func CreateCommand(car *state.Car) Command {
	car.State().Changes()

	return Command{
		Id: car.Id(),
		Changes: changes{
			Race: car.Race().State().Changes(),
			Car:  car.State().Changes(),
		},
	}
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
