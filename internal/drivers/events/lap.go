package events

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"time"
)

type Lap interface {
	drivers.Event
	Number() uint32
	Time() time.Duration
	Recorded() time.Duration
}

type GenericLap struct {
	drivers.Event
	number   uint32
	time     time.Duration
	recorded time.Duration
}

func NewLap(car drivers.Car, number uint32, time time.Duration, recorded time.Duration) Lap {
	return &GenericLap{Event: NewGenericEvent(car), number: number, time: time, recorded: recorded}
}

func (g GenericLap) Number() uint32 {
	return g.number
}

func (g GenericLap) Time() time.Duration {
	return g.time
}

func (g GenericLap) Recorded() time.Duration {
	return g.recorded
}
