package implement

import "github.com/qvistgaard/openrms/internal/state"

type Implementer interface {
	EventLoop() error
	EventChannel() <-chan Event
	SendCarState(c state.CarChanges) error
	SendRaceState(r state.CourseChanges) error

	// Resend relevant car state to implement.
	//
	// this method is executed if for example the controller
	// looses link with the dongle. But also for each car if
	// race status changes.
	ResendCarState(c *state.Car)
}
