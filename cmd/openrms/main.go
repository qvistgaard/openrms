package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/state"
	"io/ioutil"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	flagConfig := flag.String("config", "../config.yaml", "OpenRMS Config file")
	file, err := ioutil.ReadFile(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	implement, err := config.CreateImplementFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	telemetry, err := config.CreateTelemetryReceiverFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := config.CreateRaceRulesFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	r := state.CreateRace(map[string]interface{}{})

	// Todo: add repository
	// todo: automatically initialize new cars if detected
	cars := map[uint8]*state.Car{
		2: state.CreateCar(r, 2, map[string]interface{}{}, rules),
	}

	wg.Add(1)
	go eventloop(implement)

	wg.Add(1)
	go processEvents(implement, telemetry, cars)

	wg.Add(1)
	go processTelemetry(telemetry)

	wg.Wait()
}
