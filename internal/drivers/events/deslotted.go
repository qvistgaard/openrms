package events

import "github.com/qvistgaard/openrms/internal/drivers"

type Deslotted interface {
	drivers.Event
	Deslotted() bool
}

type GenericDeslotted struct {
	drivers.Event
	deslotted bool
}

// NewDeslotted
// Deprecated: use on track instead
func NewDeslotted(car drivers.Car, deslotted bool) Deslotted {
	return GenericDeslotted{NewGenericEvent(car), deslotted}
}

func (g GenericDeslotted) Deslotted() bool {
	return g.deslotted
}
