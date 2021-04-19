package implement

import "github.com/qvistgaard/openrms/internal/state"

type Implementer interface {
	EventLoop() error
	EventChannel() <-chan Event
	SendCarState(c state.CarChanges) error
	SendRaceState(r state.CourseChanges) error
}
