package main

import (
	"log"
	"openrms/implement"
	"openrms/state"
	"openrms/telemetry"
)

func eventloop(i implement.Implementer) error {
	defer wg.Done()
	log.Println("started race openrms.connector.")
	err := i.EventLoop()
	log.Println(err)
	return err
}

func processEvents(i implement.Implementer, telemetry telemetry.Receiver, cars map[uint8]*state.Car) {
	defer wg.Done()

	log.Println("started event processor.")
	for {
		e, _ := i.WaitForEvent()
		c := cars[e.Id]
		if c != nil {
			c.State().Get(state.CarEvent).Set(e)
			i.SendCommand(implement.CreateCommand(c))
			telemetry.CarChanges(c)
			c.State().ResetChanges()
		}
	}
}

func processTelemetry(receiver telemetry.Receiver) {
	defer wg.Done()
	if receiver != nil {
		log.Println("started telemetry receiver.")
		receiver.Process()
	} else {
		log.Println("telemetry receiver disabled.")
	}
}
