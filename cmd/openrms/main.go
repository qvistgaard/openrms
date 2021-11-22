package main

import (
	"flag"
	"github.com/pkg/browser"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/rms"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func main() {
	var err error

	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	flagLogfile := flag.String("log-file", "openrms.log", "OpenRMS log file")
	flagLoglevel := flag.String("log-level", "info", "Log level")
	flagBrowser := flag.Bool("open-browser", true, "Open browser on launch")
	flag.Parse()

	level, err := log.ParseLevel(*flagLoglevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.SetReportCaller(false)

	logFile, err := os.OpenFile(*flagLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	c := &application.Context{}
	err = config.FromFile(c, flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateImplement(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateWebserver(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreatePostProcessors(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateRules(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateCarRepository(c)
	if err != nil {
		log.Fatal(err)
	}
	err = config.ConfigureTrack(c)
	if err != nil {
		log.Fatal(err)
	}

	err = config.ConfigureRace(c)
	if err != nil {
		log.Fatal(err)
	}

	if *flagBrowser {
		browser.OpenURL("http://localhost:8080")
	}

	rms.Create(c).Run()

}
