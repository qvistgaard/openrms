package main

import (
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/tarm/serial"
	"log"
	"openrms/plugins/connector/oxigen"
	"openrms/plugins/rules/fuel"
	"openrms/state"

	"openrms/plugins/telemetry/influxdb"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var err error

	c := &serial.Config{Name: "COM5", Baud: 19200, ReadTimeout: time.Millisecond * 100}
	connection, _ := serial.OpenPort(c)
	o, err := oxigen.Connect(connection)
	if err != nil {
		log.Fatal(err)
	}

	input := queue.NewFIFO()
	output := queue.NewFIFO()
	telemetry := queue.NewFIFO()

	// Replace with dynamic list
	targets := []state.Rule{
		new(fuel.Consumption),
	}

	// Todo add repository
	cars := map[uint8]*state.Car{
		2: state.CreateCar(1, map[string]interface{}{}, targets),
	}

	wg.Add(1)
	// todo: change args to struct
	go eventloop(o, input, output)
	wg.Add(1)
	// todo: change args to struct
	go processEvents(o, output, telemetry, targets, cars)
	wg.Add(1)
	go processTelemetry(influxdb.Connect(), telemetry)

	wg.Wait()
}
