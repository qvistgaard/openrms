package main

import (
	"flag"
	"github.com/madflojo/tasks"
	"github.com/pkg/browser"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/rms"
	"github.com/qvistgaard/openrms/internal/tui"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	var err error
	rand.Seed(time.Now().UnixNano())

	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	flagLogfile := flag.String("log-file", "openrms.log", "OpenRMS log file")
	flagLoglevel := flag.String("log-level", "debug", "Log level")
	flagBrowser := flag.Bool("open-browser", false, "Open browser on launch")
	flagImplement := flag.String("driver", "", "Driver")
	flag.Parse()

	level, err := log.ParseLevel(*flagLoglevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.SetReportCaller(false)

	logFile, err := os.OpenFile(*flagLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	mw := io.MultiWriter(logFile)
	log.SetOutput(mw)

	c := &application.Context{
		Scheduler: tasks.New(),
	}
	defer c.Scheduler.Stop()

	err = config.FromFile(c, flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = config.CreateImplement(c, flagImplement)
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
	err = config.CreateValueFactory(c)
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

	b := tui.CreateBridge(c.Leaderboard, c.Scheduler, c.Cars, c.Race)

	if *flagBrowser {
		browser.OpenURL("http://localhost:8080")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go rms.Create(&wg, c.Postprocessors, c.Implement, c.Track, c.Rules, c.Race, c.Cars).Run()
	//go c.Webserver.RunServer(&wg)

	b.Run()
	b.UI.Run()
}
