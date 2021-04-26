package rms

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Runner struct {
	context *context.Context
	wg      sync.WaitGroup
}

func (r *Runner) Run() {
	r.wg.Add(1)
	go r.eventloop()

	r.wg.Add(1)
	go r.processEvents()

	r.wg.Wait()
}

func Create(c *context.Context) *Runner {
	return &Runner{context: c}
}

func (r *Runner) eventloop() error {
	defer func() {
		log.Fatal("Eventloop died")
	}()
	defer r.wg.Done()
	log.Info("started race OpenRMS connector.")
	err := r.context.Implement.EventLoop()
	log.Println(err)
	return err
}

func (r *Runner) processEvents() {
	defer func() {
		log.Fatal("process events died")
	}()
	defer r.wg.Done()

	log.Info("started event processor.")

	go r.processCommands()
	for {
		select {
		case e := <-r.context.Implement.EventChannel():
			if e.Id > 0 {
				if c, ok := r.context.Cars.Get(e.Id); ok {
					e.SetCarState(c)
					carChanges := c.Changes()
					if len(carChanges.Changes) > 0 {
						r.context.Implement.SendCarState(carChanges)
						r.context.Postprocessors.PostProcessCar(c.Changes())
					}
					c.ResetStateChangeStatus()
				}
			}
			raceChanges := r.context.Course.Changes()
			if len(raceChanges.Changes) > 0 {
				r.context.Implement.SendRaceState(raceChanges)
				r.context.Postprocessors.PostProcessRace(raceChanges)
			}
			r.context.Course.ResetStateChangeStatus()
		}
	}
}

func (r *Runner) processCommands() {
	for {
		select {
		case command := <-r.context.Postprocessors.CommandChannel:
			log.Infof("Received command: %T, %+v", command, command)
			if cc, ok := command.(state.CarCommand); ok {
				if c, ok := r.context.Cars.Get(cc.CarId); ok {
					c.Set(cc.Name, cc.Value)
				}
			} else if cc, ok := command.(state.CourseCommand); ok {
				r.context.Course.Set(cc.Name, cc.Value)
			}
		}
	}
}
