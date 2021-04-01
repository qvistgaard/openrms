package telemetry

import (
	"openrms/state"
	"time"
)
import queue "github.com/enriquebris/goconcurrentqueue"

type Car struct {
	Car    uint8
	Values map[string]interface{}
	Time   time.Time
}

type Processor interface {
	ProcessCar(t Car)
	// ProcessRace(t Car)
}

type Receiver interface {
}

type QueueReceiver struct {
	queue queue.Queue
}

func (q *QueueReceiver) CarChanges(values map[string]state.StateInterface) {
	for k, v := range values {
		car, err := v.Owner().(state.Car)
		if !err {
			q.queue.Enqueue(&Car{
				Car: car.Id(),
			})
		}

	}

	q.queue.Enqueue()
}
