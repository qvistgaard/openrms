package pit

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

const (
	State   = "pit-rule-pit-stop-state"
	Started = "started"
	Stopped = "stopped"
	Exiting = "exiting"
	Locked  = "locked"
)

type Pit struct {
	rules state.Rules
}

func CreatePitRule(ctx *context.Context) state.Rule {
	p := new(Pit)
	p.rules = ctx.Rules
	return p
}

func (p *Pit) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if v.Name() == state.ControllerTriggerValue {
			triggerValue := v.Get().(state.TriggerValue)
			if c.Get(state.CarInPit).(bool) {
				rules := p.rules.PitRules()
				if len(rules) > 0 {
					cancel := make(chan bool)
					o := make(chan state.PitRule, len(rules))
					if triggerValue == 0 && c.Get(State) == Stopped {
						c.Set(State, Started)
						for _, pr := range rules {
							o <- pr
						}
						go p.handlePitStop(c, o, cancel)
					} else if c.Get(State) != Locked {
						c.Set(State, Exiting)
						cancel <- true
						close(cancel)
						close(o)
					}
				}
			} else {
				c.Set(State, Stopped)
			}
		}
	}
}

func (p *Pit) InitializeCarState(c *state.Car) {
	c.Set(State, Stopped)
	c.Subscribe(state.ControllerTriggerValue, p)
}

func (p *Pit) InitializeCourseState(race *state.Course) {
	//panic("implement me")
}

func (p *Pit) handlePitStop(c *state.Car, o chan state.PitRule, cancel chan bool) {
	select {
	case r := <-o:
		r.HandlePitStop(c, cancel)
	}
}
