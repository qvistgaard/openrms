package fuel

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefuellingPitSequence(t *testing.T) {
	carState := &state{}
	carState.config = FuelConfig{
		TankSize: 10,
		FlowRate: 1,
	}
	carState.fuel = observable.Create(float32(carState.config.TankSize))

	createFuelObserver(200, fuelModifier(carState), func(f float32) {})
	carState.consumed = 10

	sequence := NewSequence(carState)

	err := sequence.Start()

	assert.Nil(t, err)
	assert.Equal(t, carState.consumed, float32(0))

}
