package events

import "github.com/qvistgaard/openrms/internal/drivers"

type OnTrack interface {
	drivers.Event
	OnTrack() bool
}

type GenericOnTrack struct {
	drivers.Event
	ontrack bool
}

func NewOnTrack(car drivers.Car, ontrack bool) OnTrack {
	return GenericOnTrack{NewGenericEvent(car), ontrack}
}

func (g GenericOnTrack) OnTrack() bool {
	return g.ontrack
}
