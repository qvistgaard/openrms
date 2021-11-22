package rms

import (
	"context"
	"github.com/qvistgaard/openrms/internal/config/application"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Runner struct {
	context *application.Context
	wg      sync.WaitGroup
}

func (r *Runner) Run() {
	r.wg.Add(1)
	go r.context.Webserver.RunServer()

	r.wg.Add(1)
	go r.eventLoop()

	r.wg.Add(1)
	go r.processEvents()

	r.wg.Wait()
}

func Create(c *application.Context) *Runner {
	return &Runner{context: c}
}

func (r *Runner) eventLoop() error {
	defer func() {
		log.Fatal("rms: Eventloop died")
	}()
	defer r.wg.Done()
	log.Info("rms: started race OpenRMS connector.")
	err := r.context.Implement.EventLoop()
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
	r.context.Postprocessors.Init(background)
	r.context.Track.Init(background, r.context.Postprocessors.ValuePostProcessor())

	for _, rule := range r.context.Rules.RaceRules() {
		rule.ConfigureRaceState(r.context.Race)
	}
	for _, rule := range r.context.Rules.RaceRules() {
		rule.InitializeRaceState(r.context.Race, background, r.context.Postprocessors.ValuePostProcessor())
	}

	r.context.Race.Init(background, r.context.Postprocessors.ValuePostProcessor())

	channel := r.context.Implement.EventChannel()
	for {
		select {
		case e := <-channel:
			start := time.Now()
			if e.Car.Id > 0 {
				if c, ok, _ := r.context.Cars.Get(e.Car.Id, background); ok {
					c.UpdateFromEvent(e)
				}
			}
			r.context.Race.UpdateTime()
			log.Tracef("processing time: %s", time.Now().Sub(start))
		}
	}
}
