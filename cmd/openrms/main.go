package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/rms"
	log "github.com/sirupsen/logrus"
)

func main() {
	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")

	// TODO: Make configurable
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)

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

	rms.Create(c).Run()

}
