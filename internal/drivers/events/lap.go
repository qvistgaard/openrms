package events

import (
	"time"
)

type Lap interface {
	Event
	Number() uint16
	Time() time.Duration
	Recorded() time.Duration
}

type GenericLap struct {
	GenericEvent
}

func (g GenericLap) Number() uint16 {
	//TODO implement me
	panic("implement me")
}

func (g GenericLap) Time() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (g GenericLap) Recorded() time.Duration {
	//TODO implement me
	panic("implement me")
}
