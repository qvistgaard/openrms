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

func CreatePitRule(ctx *context.Context) *Pit {
	p := new(Pit)
	p.rules = ctx.Rules
	p.stops = make(map[state.CarId]chan bool)
	return p
}

func (p *Pit) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if v.Name() == state.CarInPit {
			if b, ok := v.Get().(bool); ok && !b {
				c.Set(State, Stopped)
				log.Debugf("EXIT PIT %+v", c.Get(state.CarInPit))
			} else {
				log.Debugf("Enter PIT %+v", c.Get(state.CarInPit))
			}
			return
		}

		if v.Name() == state.ControllerTriggerValue {
			if b, ok := c.Get(state.CarInPit).(bool); ok && b {
				triggerValue := v.Get().(state.TriggerValue)
				if triggerValue == 0 && c.Get(State) == Stopped {
					log.Debugf("Start PIT %+v", c.Get(state.CarInPit))

					c.Set(State, Started)
				} else if triggerValue > 0 && c.Get(State) == Started {
					log.Debugf("Cancel PIT %+v", c.Get(state.CarInPit))
					c.Set(State, Cancelled)
				}
			}
			return
		}

		if v.Name() == State {
			if v.Get().(string) == Started {
				log.Debugf("Start pithandler %+v", c.Get(state.CarInPit))
				go p.handlePitStop(c, p.stops[c.Id()])
			} else if v.Get().(string) == Cancelled {
				log.Debugf("Cancelled PIT %+v", c.Get(state.CarInPit))
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
	log.Debugf("Pit stop handler startet")
	defer func() {
		for len(cancel) > 0 {
			log.Debugf("FLushing cancel")
			<-cancel
		}
		log.Debugf("Pit stop handler ended")
	}()
	for _, r := range p.rules.PitRules() {
		r.HandlePitStop(c, cancel)
	}
}
