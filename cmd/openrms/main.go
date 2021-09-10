package main

import (
	"flag"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/rms"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func main() {
	var err error

	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	flagLogfile := flag.String("log-file", "openrms.log", "OpenRMS log file")

	// TODO: Make configurable
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(false)

	logFile, err := os.OpenFile(*flagLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	c := &context.Context{}
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
