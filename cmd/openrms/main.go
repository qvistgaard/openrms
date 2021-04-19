package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/postprocess"
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

	course, err := config.CreateCourseFromConfig(file, rules)
	if err != nil {
		log.Fatal(err)
	}

	processors, err := config.CreatePostProcessors(file)
	if err != nil {
		log.Fatal(err)
	}
	postProcess := postprocess.CreatePostProcess(processors)

	wg.Add(1)
	go eventloop(implement)

	wg.Add(1)
	go processEvents(implement, postProcess, repository, course, rules)

	wg.Wait()
}
