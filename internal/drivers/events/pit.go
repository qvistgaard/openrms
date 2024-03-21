package events

import "github.com/qvistgaard/openrms/internal/drivers"

type InPit interface {
	drivers.Event
	InPit() bool
}

type GenericInPit struct {
	drivers.Event
	inPit bool
}

func NewInPit(car drivers.Car, inPit bool) InPit {
	return GenericInPit{NewGenericEvent(car), inPit}
}

func (g GenericInPit) InPit() bool {
	return g.inPit
}
