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
	Exiting   = "exiting"
	Locked    = "locked"
)

type Pit struct {
	rules state.Rules
	stops map[state.CarId]chan bool
}

func CreatePitRule(ctx *context.Context) state.Rule {
	p := new(Pit)
	p.rules = ctx.Rules
	p.stops = make(map[state.CarId]chan bool)
	return p
}

func (p *Pit) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		log.Infof("%s, %+v, %+v, %+v", v.Name(), c.Get(state.CarInPit), c.Get(State), c.Get(state.ControllerTriggerValue))
		if v.Name() == state.CarInPit {
			if b, ok := v.Get().(bool); ok && !b {
				c.Set(State, Stopped)
				log.Infof("EXIT PIT %+v", c.Get(state.CarInPit))
			} else {
				log.Infof("Enter PIT %+v", c.Get(state.CarInPit))
			}
			return
		}

		if v.Name() == state.ControllerTriggerValue {
			if b, ok := c.Get(state.CarInPit).(bool); ok && b {
				triggerValue := v.Get().(state.TriggerValue)
				if triggerValue == 0 && c.Get(State) == Stopped {
					log.Infof("Start PIT %+v", c.Get(state.CarInPit))

					c.Set(State, Started)
				} else if triggerValue > 0 && c.Get(State) == Started {
					log.Infof("Cancel PIT %+v", c.Get(state.CarInPit))
					c.Set(State, Cancelled)
				}
			}
			return
		}

		if v.Name() == State {
			if v.Get().(string) == Started {
				log.Infof("Run PIT %+v", c.Get(state.CarInPit))
				go p.handlePitStop(c, p.stops[c.Id()])
			} else if v.Get().(string) == Cancelled {
				log.Infof("Cancelled PIT %+v", c.Get(state.CarInPit))
				p.stops[c.Id()] <- true
			}
			return
		}
	}
}

func (p *Pit) InitializeCarState(c *state.Car) {
	p.stops[c.Id()] = make(chan bool, 10)
	c.Set(State, Stopped)
	c.Subscribe(state.ControllerTriggerValue, p)
	c.Subscribe(state.CarInPit, p)
	c.Subscribe(State, p)
}

func (p *Pit) InitializeCourseState(race *state.Course) {
	//panic("implement me")
}

func (p *Pit) handlePitStop(c *state.Car, cancel chan bool) {
	defer func() {
		for len(cancel) > 0 {
			log.Info("FLushing cancel")
			<-cancel
		}
		log.Info("Pit stop handler ended")
	}()
	for _, r := range p.rules.PitRules() {
		r.HandlePitStop(c, cancel)
	}
}
