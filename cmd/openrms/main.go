package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/config/context"
	log "github.com/sirupsen/logrus"
	"sync"
)

var wg sync.WaitGroup

func main() {
	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")

	c := &context.Context{}
	var err error
	err = config.FromFile(c, flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateImplement(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateRules(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreatePostProcessors(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateCarRepository(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateCourse(c)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Make configurable
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)

	wg.Add(1)
	go eventloop(c.Implement)

	wg.Add(1)
	// go processEvents(implement, postProcess, repository, course, rules)
	go processEvents(c.Implement, c.Postprocessors, c.Cars, c.Course, c.Rules)

	wg.Wait()

}
