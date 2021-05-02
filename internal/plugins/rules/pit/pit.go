package pit

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
)

const (
	State     = "pit-rule-pit-stop-state"
	Started   = "started"
	Stopped   = "stopped"
	Cancelled = "cancelled"
	Complete  = "complete"
	Locked    = "locked"
)

type Rule struct {
	rules  state.Rules
	stops  map[state.CarId]chan bool
	course *state.Course
}

func CreatePitRule(ctx *context.Context) *Rule {
	p := new(Rule)
	p.rules = ctx.Rules
	p.stops = make(map[state.CarId]chan bool)
	return p
}

func (p *Rule) Notify(v *state.Value) {
	if p.course.Get(state.RaceStatus) != state.RaceStatusStopped {
		if c, ok := v.Owner().(*state.Car); ok {
			if v.Name() == state.CarInPit {
				if b, ok := v.Get().(bool); ok && !b {
					c.Set(State, Stopped)
					log.WithField("car", c.Id()).
						WithField(State, v.Get()).
						Debugf("pit handler: car exited pitlane")
				} else {
					log.WithField("car", c.Id()).
						WithField(State, v.Get()).
						Debugf("pit handler: car entered pitlane")
				}
				return
			}

			if v.Name() == state.ControllerTriggerValue {
				if b, ok := c.Get(state.CarInPit).(bool); ok && b {
					triggerValue := v.Get().(state.TriggerValue)
					if triggerValue == 0 && c.Get(State) == Stopped {
						log.WithField("car", c.Id()).
							WithField("triggerValue", triggerValue).
							WithField(State, v.Get()).
							Debugf("pit handler: detected triggerValue change")
						c.Set(State, Started)
					} else if triggerValue > 0 && c.Get(State) == Started {
						log.WithField("car", c.Id()).
							WithField(State, v.Get()).
							WithField("triggerValue", triggerValue).
							Debugf("pit handler: detected triggerValue change")
						c.Set(State, Cancelled)
					}
				}
				return
			}

			if v.Name() == State {
				if v.Get().(string) == Started {
					log.WithField("car", c.Id()).
						WithField(State, v.Get()).
						Debugf("pit handler: pit stop started")
					go p.handlePitStop(c, p.stops[c.Id()])
				} else if v.Get().(string) == Locked {
					c.Set(state.CarMaxSpeed, state.Speed(0))
				} else if v.Get().(string) == Cancelled {
					log.WithField("car", c.Id()).
						WithField(State, v.Get()).
						Debugf("pit handler: pit stop cancelled")
					p.stops[c.Id()] <- true
				}
				return
			}
		}
	}
}

func (p *Rule) InitializeCarState(c *state.Car) {
	p.stops[c.Id()] = make(chan bool, 10)
	c.Set(State, Stopped)
	c.Subscribe(state.ControllerTriggerValue, p)
	c.Subscribe(state.CarInPit, p)
	c.Subscribe(State, p)
}

func (p *Rule) InitializeCourseState(course *state.Course) {
	p.course = course
}

func (p *Rule) handlePitStop(c *state.Car, cancel chan bool) {
	log.WithField("car", c.Id()).Debugf("pit handler: started")
	defer func() {
		for len(cancel) > 0 {
			log.WithField("car", c.Id()).Debugf("pit handler: flushing cancel channel")
			<-cancel
		}
		log.WithField("car", c.Id()).Debugf("pit handler: ended")
	}()
	for _, r := range p.rules.PitRules() {
		if !r.HandlePitStop(c, cancel) {
			break
		}
	}
	c.SetDefault(state.CarMaxSpeed)
	c.Set(State, Complete)
	log.WithField("car", c.Id()).
		WithField(State, c.Get(State)).
		Debugf("pit handler: ended")
}
