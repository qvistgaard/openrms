package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"io/ioutil"

	"sync"
)

var wg sync.WaitGroup

func main() {

	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	file, err := ioutil.ReadFile(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Make configurable
	log.SetLevel(log.InfoLevel)

	implement, err := config.CreateImplementFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := config.CreateRaceRulesFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := config.CreateCarRepositoryFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	// todo: Create race from configuration
	r := state.CreateRace(map[string]interface{}{}, rules)

	// Todo: add repository
	// todo: automatically initialize new cars if detected
	/*	cars := map[uint8]*state.Car{
			2: state.CreateCar(r, 2, map[string]interface{}{}, rules),
		}
	*/
	processors, err := config.CreatePostProcessors(file)
	if err != nil {
		log.Fatal(err)
	}
	postProcess := postprocess.CreatePostProcess(processors)

	wg.Add(1)
	go eventloop(implement)

	wg.Add(1)
	go processEvents(implement, postProcess, repository, r, rules)
	/*
		wg.Add(1)
		go postProcess.PostProcess()
	*/
	wg.Wait()
}
