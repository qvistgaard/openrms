package pit

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSomething(t *testing.T) {
	r := state.CreateRuleList()
	p := CreatePitRule(r)
	r.Append(p)

	c := state.CreateCar(nil, 1, nil, r)
	c.Set(fuel.CarConfigFuel, float32(100))
	c.Set(fuel.CarFuel, float32(50))
	c.Set(state.CarInPit, true)
	c.Set(state.ControllerTriggerValue, uint8(0))

	assert.Equal(t, float32(100), c.Get(fuel.CarFuel))
}
