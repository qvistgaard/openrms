package rms

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
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
	c.Course.Set(state.RMSStatus, state.Initialized)
	return &Runner{context: c}
}

func (r *Runner) eventloop() error {
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

	go r.processCommands()

	r.context.Implement.SendRaceState(r.context.Course.State())

	for {
		select {
		case e := <-r.context.Implement.EventChannel():
			start := time.Now()
			if e.Id > 0 {
				if c, ok, created := r.context.Cars.Get(e.Id); ok {
					e.SetCarState(c)
					e.SetCourseState(r.context.Course)
					if created {
						log.WithField("car", e.Id).
							Info("new car found, resending state")
						r.context.Implement.ResendCarState(c)
					} else {
						carChanges := c.Changes()
						if len(carChanges.Changes) > 0 {
							r.context.Implement.SendCarState(carChanges)
							r.context.Postprocessors.PostProcessCar(c.Changes())
						}
					}
					c.ResetStateChangeStatus()
				}
			}
			raceChanges := r.context.Course.Changes()
			if len(raceChanges.Changes) > 0 {
				if r.context.Course.IsChanged(state.RaceStatus) {
					for _, c := range r.context.Cars.All() {
						r.context.Implement.ResendCarState(c)
					}
				}
				r.context.Implement.SendRaceState(raceChanges)
				r.context.Postprocessors.PostProcessRace(raceChanges)
			}
			r.context.Course.ResetStateChangeStatus()
			log.Tracef("processing time: %s", time.Now().Sub(start))
		}
	}
}

func (r *Runner) processCommands() {
	for {
		select {
		case command := <-r.context.Postprocessors.CommandChannel:
			log.Infof("Received command: %T, %+v", command, command)
			if cc, ok := command.(state.CarCommand); ok {
				if r.context.Cars.Exists(cc.CarId) {
					c, _, _ := r.context.Cars.Get(cc.CarId)
					c.Set(cc.Name, cc.Value)
				}
			} else if cc, ok := command.(state.CourseCommand); ok {
				r.context.Course.Set(cc.Name, cc.Value)
			}
		}
	}
}
