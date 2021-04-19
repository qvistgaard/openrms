package pit

import (
	"github.com/qvistgaard/openrms/internal/state"
)

type Pit struct{}

func (p *Pit) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		if v.Name() == state.CarInPit {

		}
	}
}

func (p *Pit) InitializeCarState(c *state.Car) {
	c.Subscribe(state.CarEventSequence, p)
}

func (p *Pit) InitializeRaceState(race *state.Race) {
	panic("implement me")
}
