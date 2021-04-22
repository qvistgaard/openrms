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

	cars := make(map[state.CarId]*state.Car)
	go processCommands(cars, postProcess, course)
	for {
		select {
		case e := <-i.EventChannel():
			log.Infof("Event: %+v", e)
			var c *state.Car
			if e.Id > 0 {
				if _, ok := cars[e.Id]; !ok {
					c = state.CreateCar(course, e.Id, repository.GetCarById(e.Id), rules)
					cars[e.Id] = c
				} else {
					c = cars[e.Id]
				}
				if c != nil {
					e.SetCarState(c)

					carChanges := c.Changes()
					if len(carChanges.Changes) > 0 {
						i.SendCarState(carChanges)
						postProcess.PostProcessCar(c.Changes())
					}
					c.ResetStateChangeStatus()

				}
			}
			raceChanges := course.Changes()
			if len(raceChanges.Changes) > 0 {
				i.SendRaceState(raceChanges)
				postProcess.PostProcessRace(raceChanges)
			}
			course.ResetStateChangeStatus()

		}
	}
}

func processCommands(cars map[state.CarId]*state.Car, postProcess postprocess.PostProcess, course *state.Course) {
	for {
		select {
		case command := <-postProcess.CommandChannel:
			log.Infof("Received command: %T, %+v", command, command)
			if cc, ok := command.(state.CarCommand); ok {
				if c, ok := cars[cc.CarId]; ok {
					c.Set(cc.Name, cc.Value)
				}
			} else if cc, ok := command.(state.CourseCommand); ok {
				course.Set(cc.Name, cc.Value)
			}
		}
	}
}
