package main

import (
	"flag"
	"github.com/madflojo/tasks"
	"github.com/pkg/browser"
	"github.com/qvistgaard/openrms/cmd/openrms/configuration"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/plugins/ontrack"
	"github.com/qvistgaard/openrms/internal/plugins/telemetry"
	"github.com/qvistgaard/openrms/internal/rms"
	"github.com/qvistgaard/openrms/internal/tui"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

func main() {
	var err error
	var wg sync.WaitGroup

	/*	wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	*/
	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	flagLogfile := flag.String("log-file", "openrms.log", "OpenRMS log file")
	flagLoglevel := flag.String("log-level", "debug", "Log level")
	flagBrowser := flag.Bool("open-browser", false, "Open browser on launch")
	flagDriver := flag.String("driver", "", "Driver")
	tuiFlag := flag.Bool("tui", true, "Enable or disable tui")
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

	track, err := configuration.Track(cfg, driver)
	if err != nil {
		log.Fatal(err)
	}

	race, err := configuration.Race(cfg, driver)
	if err != nil {
		log.Fatal(err)
	}

	/*	err = config.CreateWebserver(c)
		if err != nil {
			log.Fatal(err)
		}*/

	confirmationPlugin, err := configuration.ConfirmationPlugin(cfg)
	if err != nil {
		log.Fatal(err)
	}

	racePlugin, err := configuration.RacePlugin(cfg, race, confirmationPlugin)
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
	pitPlugin, _ := configuration.PitPlugin(cfg, fuelPlugin, limpModePlugin)
	leaderboardPlugin := telemetry.New(fuelPlugin, limpModePlugin, pitPlugin)
	flagPlugin := flags.New(racePlugin)
	ontrackPlugin := ontrack.New(flagPlugin)

	plugins, err := configuration.Plugins(cfg)
	if err != nil {
		log.Fatal(err)
	}
	plugins.Append(racePlugin)
	plugins.Append(pitPlugin)
	plugins.Append(leaderboardPlugin)
	plugins.Append(limpModePlugin)
	plugins.Append(fuelPlugin)
	plugins.Append(flagPlugin)
	plugins.Append(confirmationPlugin)
	plugins.Append(ontrackPlugin)

	repository, err := configuration.CarRepository(cfg, driver, plugins)
	if err != nil {
		log.Fatal(err)
	}

	if *flagBrowser {
		err = browser.OpenURL("http://localhost:8080")
		if err != nil {
			log.Error(err)
		}
	}

	wg.Add(1)
	go rms.Create(&wg, driver, plugins, track, race, repository).Run()
	//go c.Webserver.RunServer(&wg)

	if !*tuiFlag {
		wg.Wait()
	} else {
		log.SetOutput(io.Writer(logFile))
		b := tui.CreateBridge(leaderboardPlugin, racePlugin, scheduler, track, repository, race, confirmationPlugin, flagPlugin)
		b.Run()
		b.UI.Run()
	}
}
