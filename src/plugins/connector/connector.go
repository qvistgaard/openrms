package connector

import "io"
import queue "github.com/enriquebris/goconcurrentqueue"

type Connector interface {
	io.Closer
	EventLoop(input queue.Queue, output queue.Queue) error
	Stop()
	Start()
	Pause()
	Flag(lc bool)
	MaxSpeed(speed uint8)
	PitLaneLapCount(enabled bool, entry bool)
	Shutdown()
}
