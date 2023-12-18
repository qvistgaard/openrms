package events

import (
	"github.com/qvistgaard/openrms/internal/drivers"
)

type GenericEvent struct {
	car drivers.Car
}

func NewGenericEvent(car drivers.Car) drivers.Event {
	return &GenericEvent{car: car}
}

func (g GenericEvent) Car() drivers.Car {
	return g.car
}
