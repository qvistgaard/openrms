package telemetry

import (
	"github.com/qvistgaard/openrms/internal/state"
	"time"
)
import queue "github.com/enriquebris/goconcurrentqueue"

func NewQueueReceiver(processor Processor) *QueueReceiver {
	q := new(QueueReceiver)
	q.queue = queue.NewFIFO()
	q.processor = processor
	return q
}

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
	CarChanges(car *state.Car)
	Process()
}

type QueueReceiver struct {
	queue     queue.Queue
	processor Processor
}

func (q *QueueReceiver) Process() {
	element, err := q.queue.DequeueOrWaitForNextElement()
	if err != nil {
		car, ok := element.(Car)
		if ok {
			q.processor.ProcessCar(car)
		}
	}
}

func (q *QueueReceiver) CarChanges(car *state.Car) {
	c := new(Car)
	c.Values = map[string]interface{}{}
	c.Time = time.Now()
	for k, v := range car.State().Changes() {
		c.Values[k] = v.Get()
	}
	q.queue.Enqueue(c)
}
