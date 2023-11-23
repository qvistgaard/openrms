package fuel

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/stretchr/testify/assert" // Use a testing library like testify/assert for assertions.
	"testing"
)

func TestMachine(t *testing.T) {
	// Define a function to use as the fuelUpdate callback.
	fuelUpdate := func(ctx context.Context, args ...interface{}) error {
		// Implement the fuel update logic for testing.
		return nil
	}

	ctx := context.TODO()

	// Test transitions and states.
	t.Run("TestCarOnTrackTransitionOnStateCarOnTrack", func(t *testing.T) {
		m := machine(fuelUpdate)

		err := m.Fire(triggerCarOnTrack)
		assert.Error(t, err)
	})

	t.Run("TestUpdateFuelOnStateCarOnTrack", func(t *testing.T) {
		m := machine(fuelUpdate)

		err := m.Fire(triggerCarOnTrack)
		assert.Error(t, err)

		err = m.Fire(triggerUpdateFuelLevel, uint8(2))
		assert.NoError(t, err)

	})

	t.Run("TestUpdateFuelOnStateCarDeslotted", func(t *testing.T) {
		m := machine(fuelUpdate)

		err := m.Fire(triggerCarDeslotted)
		assert.NoError(t, err)

		err = m.Fire(triggerUpdateFuelLevel, uint8(2))
		assert.Error(t, err)
	})

	t.Run("TestCarDeslottedTransition", func(t *testing.T) {
		m := machine(fuelUpdate)

		m.Fire(triggerCarDeslotted)
		err := m.Fire(triggerCarDeslotted)

		assert.NoError(t, err)
		state, err := m.State(ctx)
		assert.NoError(t, err)

		assert.Equal(t, stateCarDeslotted, state)
	})

	t.Run("TestUpdateFuelLevelTransition", func(t *testing.T) {
		m := machine(fuelUpdate)

		err := m.Fire(triggerUpdateFuelLevel, types.Percent(50))
		assert.NoError(t, err)
		state, err := m.State(ctx)
		assert.Equal(t, stateCarOnTrack, state)
	})
}
