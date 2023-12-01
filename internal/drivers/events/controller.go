package events

import "github.com/qvistgaard/openrms/internal/drivers"

type ControllerBatteryWarning interface {
	drivers.Event
	BatteryWarning() bool
}

type ControllerLinkEvent interface {
	drivers.Event
	Link() bool
}
