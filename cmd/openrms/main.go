package main

import (
	"flag"
	"github.com/madflojo/tasks"
	"github.com/pkg/browser"
	"github.com/qvistgaard/openrms/cmd/openrms/configuration"
	"github.com/qvistgaard/openrms/internal/plugins/leaderboard"
	"github.com/qvistgaard/openrms/internal/rms"
	"github.com/qvistgaard/openrms/internal/tui"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

func main() {
	var err error
	// rand.Seed(time.Now().UnixNano())

	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	flagLogfile := flag.String("log-file", "openrms.log", "OpenRMS log file")
	flagLoglevel := flag.String("log-level", "debug", "Log level")
	flagBrowser := flag.Bool("open-browser", false, "Open browser on launch")
	flagDriver := flag.String("driver", "", "Driver")
	flag.Parse()

	level, err := log.ParseLevel(*flagLoglevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.SetReportCaller(false)

	logFile, err := os.OpenFile(*flagLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	scheduler := tasks.New()
	defer scheduler.Stop()

	cfg, err := configuration.FromFile(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := configuration.Driver(cfg, flagDriver)
	if err != nil {
		log.Fatal(err)
	}

	/*	err = config.CreateWebserver(c)
		if err != nil {
			log.Fatal(err)
		}*/

	racePlugin, err := configuration.RacePlugin(cfg)
	if err != nil {
		log.Fatal(err)
	}

	limpModePlugin, err := configuration.LimbModePlugin(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fuelPlugin, err := configuration.FuelPlugin(cfg, limpModePlugin)
	if err != nil {
		log.Fatal(err)
	}
	leaderboardPlugin := leaderboard.New(fuelPlugin, limpModePlugin)

	plugins, err := configuration.Plugins(cfg)
	if err != nil {
		log.Fatal(err)
	}
	plugins.Append(racePlugin)
	plugins.Append(leaderboardPlugin)
	plugins.Append(limpModePlugin)
	plugins.Append(fuelPlugin)

	repository, err := configuration.CarRepository(cfg, driver, plugins)
	if err != nil {
		log.Fatal(err)
	}
	track, err := configuration.Track(cfg, driver)
	if err != nil {
		log.Fatal(err)
	}

	race, err := configuration.Race(cfg, driver)
	if err != nil {
		log.Fatal(err)
	}

	if *flagBrowser {
		browser.OpenURL("http://localhost:8080")
	}

	b := tui.CreateBridge(leaderboardPlugin, racePlugin, scheduler, repository, race)

	var wg sync.WaitGroup
	wg.Add(1)
	go rms.Create(&wg, driver, plugins, track, race, repository).Run()
	//go c.Webserver.RunServer(&wg)

	// wg.Wait()
	log.SetOutput(io.Writer(logFile))

	b.Run()
	b.UI.Run()

}
