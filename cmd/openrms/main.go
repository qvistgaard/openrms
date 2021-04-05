package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"io/ioutil"

	"sync"
)

var wg sync.WaitGroup

func main() {
	log.SetLevel(log.TraceLevel)
	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	file, err := ioutil.ReadFile(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	// log.SetReportCaller(true)

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

	// todo: Create race from configuration
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
