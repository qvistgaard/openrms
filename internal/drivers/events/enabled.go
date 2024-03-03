package events

import "github.com/qvistgaard/openrms/internal/drivers"

type Enabled interface {
	drivers.Event
	Enabled() bool
}

type GenericEnabled struct {
	drivers.Event
	enabled bool
}

func NewEnabled(car drivers.Car, enabled bool) Enabled {
	return &GenericEnabled{NewGenericEvent(car), enabled}
}

func (g GenericEnabled) Enabled() bool {
	return g.enabled
}
