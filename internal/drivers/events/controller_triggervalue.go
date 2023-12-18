package events

import "github.com/qvistgaard/openrms/internal/drivers"

type ControllerTriggerValueEvent interface {
	TriggerValue() float64
}

type GenericControllerTriggerValueEvent struct {
	drivers.Event
	triggerValue float64
}

func NewControllerTriggerValueEvent(car drivers.Car, triggerValue float64) drivers.Event {
	return &GenericControllerTriggerValueEvent{Event: NewGenericEvent(car), triggerValue: triggerValue}
}

func (g GenericControllerTriggerValueEvent) TriggerValue() float64 {
	return g.triggerValue
}
