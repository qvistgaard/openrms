package main

import (
	queue "github.com/enriquebris/goconcurrentqueue"
	"log"
	"openrms/ipc"
	"openrms/plugins/connector"
	"openrms/state"
	"openrms/telemetry"
)

func eventloop(o connector.Connector, input queue.Queue, output queue.Queue) error {
	defer wg.Done()
	log.Println("started race openrms.connector.")
	err := o.EventLoop(input, output)
	log.Println(err)
	return err
}

func processEvents(o connector.Connector, output queue.Queue, telemetry telemetry.Receiver, rules []state.Rule, cars map[uint8]*state.Car) {
	defer wg.Done()

	log.Println("started event processor.")
	for {
		event, _ := output.DequeueOrWaitForNextElement()
		e := event.(ipc.Event)

		c := cars[e.Id]
		if c != nil {
			c.Get(state.RaceEvent).Set(e)
		}
		log.Printf("%+v", e)

		changes := c.StateChanges()
		telemetry.Enqueue(changes)
		log.Printf("%+v", changes)

		c.ResetChanges()
	}
}

func processTelemetry(processor telemetry.Receiver, input queue.Queue) {
	defer wg.Done()

	if processor != nil {
		log.Println("started telemetry processor.")
		for {
			event, _ := input.DequeueOrWaitForNextElement()
			t := event.(telemetry.Telemetry)
			processor.Process(t)
		}
	} else {
		log.Println("telemetry processor disabled.")
	}
}
