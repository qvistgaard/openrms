package connector

import "io"
import queue "github.com/enriquebris/goconcurrentqueue"

type RaceEvent struct {
}

type Connector interface {
	io.Closer
	EventLoop(input queue.Queue, output queue.Queue)
	stop() bool
	start() bool
	pause() bool
	flag(lc bool)
	maxSpeed(speed uint8)
	pitLaneLapCount(enabled bool, entry bool)
}
