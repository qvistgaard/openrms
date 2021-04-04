package implement

type Implementer interface {
	EventLoop() error
	WaitForEvent() (Event, error)
	SendCommand(c Command) error
}
