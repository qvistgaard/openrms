package pit

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleTestPitRule struct{}

func (SimpleTestPitRule) HandlePitStop(car *state.Car, cancel chan bool) {
	car.Set(fuel.CarFuel, float32(100))
}

func (SimpleTestPitRule) Priority() uint8 {
	// panic("implement me")
	return 50
}

func (r SimpleTestPitRule) InitializeCarState(car *state.Car) {
}

func (r SimpleTestPitRule) InitializeCourseState(race *state.Course) {
}

func TestSomething(t *testing.T) {
	ctx := &context.Context{
		Rules: state.CreateRuleList(),
	}
	p := CreatePitRule(ctx)
	ctx.Rules.Append(p)
	ctx.Rules.Append(SimpleTestPitRule{})

	c := state.CreateCar(1, nil, ctx.Rules)
	c.Set(fuel.CarConfigFuel, float32(100))
	c.Set(fuel.CarFuel, float32(50))
	c.Set(state.CarInPit, true)
	c.Set(state.ControllerTriggerValue, state.TriggerValue(0))

	assert.Equal(t, float32(100), c.Get(fuel.CarFuel))
}

func TestChannel(t *testing.T) {
	c := make(chan bool, 1)
	close(c)
	c <- true

}
