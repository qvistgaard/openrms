package ipc

import "time"

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
