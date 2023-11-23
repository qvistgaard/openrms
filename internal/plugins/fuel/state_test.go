package fuel

import (
	"context"
	"github.com/stretchr/testify/assert" // Use a testing library like testify/assert for assertions.
	"testing"
)

func TestMachine(t *testing.T) {
	// Define a function to use as the fuelUpdate callback.

	ctx := context.TODO()

	// Test transitions and states.
	t.Run("TestCarOnTrackTransitionOnStateCarOnTrack", func(t *testing.T) {
		m := machine(nil)

		err := m.Fire(triggerCarOnTrack)
		assert.Error(t, err)
	})

	t.Run("TestUpdateFuelOnStateCarOnTrack", func(t *testing.T) {
		fuel := uint8(10)
		fuelUpdate := func(ctx context.Context, args ...interface{}) error {
			// Implement the fuel update logic for testing.
			fuel, _ = args[0].(uint8)
			return nil
		}

		m := machine(fuelUpdate)

		err := m.Fire(triggerCarOnTrack)
		assert.Error(t, err)

		err = m.Fire(triggerUpdateFuelLevel, uint8(2))
		assert.NoError(t, err)
		assert.Equal(t, uint8(2), fuel)

		state, err := m.State(ctx)
		assert.NoError(t, err)

		assert.Equal(t, stateCarOnTrack, state)
	})

	t.Run("TestUpdateFuelOnStateCarDeslotted", func(t *testing.T) {
		m := machine(nil)

		err := m.Fire(triggerCarDeslotted)
		assert.NoError(t, err)

		err = m.Fire(triggerUpdateFuelLevel, uint8(2))
		assert.Error(t, err)

	})

	t.Run("TestCarDeslottedTransition", func(t *testing.T) {
		m := machine(nil)

		m.Fire(triggerCarDeslotted)
		err := m.Fire(triggerCarDeslotted)

		assert.Error(t, err)
		state, err := m.State(ctx)
		assert.NoError(t, err)

		assert.Equal(t, stateCarDeslotted, state)
	})

	t.Run("TestTransitionFromDeslottedToOnTrack", func(t *testing.T) {
		m := machine(nil)

		err := m.Fire(triggerCarDeslotted)
		assert.NoError(t, err)

		err = m.Fire(triggerCarOnTrack)
		assert.NoError(t, err)

	})
}
