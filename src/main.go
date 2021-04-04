package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	config2 "openrms/config"
	"openrms/plugins/rules/fuel"
	"openrms/state"
	"openrms/telemetry"
	"sync"
)

var wg sync.WaitGroup

type RaceContainer struct {
	race struct {
		maxSpeed uint8 `yaml:"max-speed"`
	}
}

func main() {
	config := flag.String("config", "config.yaml", "OpenRMS Config file")
	file, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", file)

	implement, err := config2.CreateImplementFromConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	t := telemetry.NewQueueReceiver(nil)
	r := state.CreateRace(map[string]interface{}{})

	// todo: Replace with dynamic list
	targets := []state.Rule{
		new(fuel.Consumption),
	}

	// Todo: add repository
	cars := map[uint8]*state.Car{
		2: state.CreateCar(r, 2, map[string]interface{}{}, targets),
	}

	wg.Add(1)
	go eventloop(implement)

	wg.Add(1)
	go processEvents(implement, t, cars)

	wg.Add(1)
	go processTelemetry(t)

	wg.Wait()
}
