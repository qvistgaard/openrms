package rms

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/plugins"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/state/car/repository"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/track"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Runner struct {
	// context        *application.Context
	wg             *sync.WaitGroup
	postprocessors *postprocess.PostProcess
	implement      implement.Implementer
	track          *track.Track
	race           *race.Race
	cars           repository.Repository
	plugins        *plugins.Plugins
}

func (r *Runner) Run() {
	r.wg.Add(1)
	go r.eventLoop()

	r.wg.Add(1)
	r.processEvents()

}

func Create(waitGroup *sync.WaitGroup, implement implement.Implementer, plugins *plugins.Plugins, track *track.Track, race *race.Race, cars repository.Repository) *Runner {
	return &Runner{
		wg:        waitGroup,
		implement: implement,
		track:     track,
		plugins:   plugins,
		race:      race,
		cars:      cars,
	}
}

func (r *Runner) eventLoop() error {
	defer func() {
		r.wg.Done()
		log.Fatal("rms: Eventloop died")
	}()
	log.Info("rms: started race OpenRMS connector.")
	err := r.implement.EventLoop()
	log.Println(err)
	return err
}

func (r *Runner) processEvents() {
	defer func() {
		panic("rms: process events died")
	}()
	defer r.wg.Done()

	log.Info("rms: started event processor.")

	background := context.Background()
	// r.postprocessors.Init(background)
	r.track.Init(background)

	for _, rule := range r.plugins.Race() {
		rule.ConfigureRace(r.race)
	}
	/*	for _, rule := range r.plugins.Race() {
		rule.InitializeRaceState(r.race, background)
	}*/

	// r.race.Init(background)

	channel := r.implement.EventChannel()
	for {
		select {
		case e := <-channel:
			start := time.Now()
			if e.Car.Id > 0 {
				if c, ok, _ := r.cars.Get(e.Car.Id); ok {
					c.UpdateFromEvent(e)
				}
			}
			r.race.UpdateFromEvent(e)
			log.Tracef("processing time: %s", time.Now().Sub(start))
		}
	}
}
