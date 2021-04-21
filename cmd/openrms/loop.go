package main

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
)

func eventloop(i implement.Implementer) error {
	defer wg.Done()
	log.Info("started race OpenRMS connector.")
	err := i.EventLoop()
	log.Println(err)
	return err
}

func processEvents(i implement.Implementer, postProcess postprocess.PostProcess, repository car.Repository, course *state.Course, rules state.Rules) {
	defer wg.Done()

	log.Info("started event processor.")

	cars := make(map[uint8]*state.Car)
	course.Set(state.RaceStatus, state.RaceStatusRunning)
	for {
		select {
		case e := <-i.EventChannel():
			var c *state.Car
			if _, ok := cars[e.Id]; !ok {
				c = state.CreateCar(course, e.Id, repository.GetCarById(e.Id), rules)
				cars[e.Id] = c
			} else {
				c = cars[e.Id]
			}
			if c != nil {
				e.SetCarState(c)

				carChanges := c.Changes()
				raceChanges := course.Changes()

				if len(raceChanges.Changes) > 0 {
					log.Infof("RACE CHANGES: %+v", raceChanges)
					i.SendRaceState(raceChanges)
					postProcess.PostProcessRace(raceChanges)
				}
				if len(carChanges.Changes) > 0 {
					i.SendCarState(carChanges)
					postProcess.PostProcessCar(c.Changes())
				}
				c.ResetStateChangeStatus()
				course.ResetStateChangeStatus()
			}
		}
	}
}
