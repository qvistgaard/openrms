package pit

import (
	"github.com/qvistgaard/openrms/internal/state"
)

type Pit struct {
	rules state.Rules
}

func CreatePitRule(rules state.Rules) state.Rule {
	p := new(Pit)
	p.rules = rules
	return p
}

// Notify TODO: Figure out how to handle set time for each supported plugin
func (p *Pit) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if v.Name() == state.ControllerTriggerValue {
			if c.Get(state.CarInPit).(bool) && v.Get().(uint8) == 0 {
				/**
				for _, pr := range p.rules.PitRules() {
					pr.HandlePitStop(c)
				}
				*/
			}
		}
	}
}

func (p *Pit) InitializeCarState(c *state.Car) {
	c.Subscribe(state.ControllerTriggerValue, p)
}

func (p *Pit) InitializeRaceState(race *state.Course) {
	//panic("implement me")
}
