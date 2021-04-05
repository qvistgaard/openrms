package main

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/telemetry"
	log "github.com/sirupsen/logrus"
)

func eventloop(i implement.Implementer) error {
	defer wg.Done()
	log.Info("started race OpenRMS connector.")
	err := i.EventLoop()
	log.Println(err)
	return err
}

func processEvents(i implement.Implementer, telemetry telemetry.Receiver, cars map[uint8]*state.Car) {
	defer wg.Done()

	log.Info("started event processor.")
	for {
		e, _ := i.WaitForEvent()
		c := cars[e.Id]
		log.WithFields(map[string]interface{}{
			"id":         e.Id,
			"ontrack":    e.Ontrack,
			"link":       e.Controller.Link,
			"in-put":     e.Car.InPit,
			"lap-number": e.LapNumber,
			"lap-time":   e.LapTime,
		}).Debug("State changed received from implement.")
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
		log.Info("started telemetry receiver.")
		receiver.Process()
	} else {
		log.Info("telemetry receiver disabled.")
	}
}
