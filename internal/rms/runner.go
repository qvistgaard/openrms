package rms

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/rules"
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
	rules          rules.Rules
	race           *race.Race
	cars           car.Repository
}

func (r *Runner) Run() {
	r.wg.Add(1)
	go r.eventLoop()

	r.wg.Add(1)
	r.processEvents()

}

func Create(waitGroup *sync.WaitGroup, postprocessors *postprocess.PostProcess, implement implement.Implementer, track *track.Track, rules rules.Rules, race *race.Race, cars car.Repository) *Runner {
	return &Runner{
		wg:             waitGroup,
		postprocessors: postprocessors,
		implement:      implement,
		track:          track,
		rules:          rules,
		race:           race,
		cars:           cars,
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
	r.postprocessors.Init(background)
	r.track.Init(background, r.postprocessors.ValuePostProcessor())

	for _, rule := range r.rules.RaceRules() {
		rule.ConfigureRaceState(r.race)
	}
	for _, rule := range r.rules.RaceRules() {
		rule.InitializeRaceState(r.race, background, r.postprocessors.ValuePostProcessor())
	}

	r.race.Init(background, r.postprocessors.ValuePostProcessor())

	channel := r.implement.EventChannel()
	for {
		select {
		case e := <-channel:
			start := time.Now()
			if e.Car.Id > 0 {
				if c, ok, _ := r.cars.Get(e.Car.Id, background); ok {
					c.UpdateFromEvent(e)
				}
			}
			r.race.UpdateTime()
			log.Tracef("processing time: %s", time.Now().Sub(start))
		}
	}
}
