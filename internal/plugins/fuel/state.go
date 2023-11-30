package fuel

import (
	"context"
	"github.com/qmuntal/stateless"
	"reflect"
)

// machineState represents the possible states of the fuel management state machine.
type machineState string

// machineTrigger represents the triggers that can cause transitions in the fuel management state machine.
type machineTrigger string

const (
	// stateCarOnTrack represents the state when a car is on the track.
	stateCarOnTrack = machineState("car-on-track")

	// stateCarDeslotted represents the state when a car is deslotted.
	stateCarDeslotted = machineState("car-deslotted")

	// triggerCarOnTrack is the trigger for transitioning a car from deslotted to on track.
	triggerCarOnTrack = machineTrigger("trigger-car-on-track")

	// triggerCarDeslotted is the trigger for transitioning a car from on track to deslotted.
	triggerCarDeslotted = machineTrigger("trigger-car-deslotted")

	// triggerUpdateFuelLevel is the trigger for updating the car's fuel level.
	triggerUpdateFuelLevel = machineTrigger("trigger-update-fuel-level")
)

// machine creates and configures a state machine for managing fuel updates.
// This state machine controls when and how the fuel level is updated based on trigger events.
// It also transitions the car's fuel state between on track and deslotted.
// The machine handles events triggered by changes in the car's trigger value and deslotting status.
func machine(fuelUpdate func(ctx context.Context, args ...any) error) *stateless.StateMachine {
	m := stateless.NewStateMachineWithMode(stateCarOnTrack, stateless.FiringImmediate)
	m.SetTriggerParameters(triggerUpdateFuelLevel, reflect.TypeOf(uint8(0)))
	m.Configure(stateCarOnTrack).
		InternalTransition(triggerUpdateFuelLevel, fuelUpdate).
		Permit(triggerCarDeslotted, stateCarDeslotted)

	m.Configure(stateCarDeslotted).
		Permit(triggerCarOnTrack, stateCarOnTrack)
	return m
}

// handleUpdateFuelLevel returns a function that updates the car's fuel level based on trigger events.
// It calculates the fuel consumption and updates the fuel level accordingly.
// The function is used as an internal transition in the state machine.
func handleUpdateFuelLevel(carState *state, size uint8, rate float32) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...interface{}) error {
		// trigger percentage
		percent := args[0].(uint8)
		if percent > 0 {
			liter := carState.fuel.Get()
			if liter > 0 {
				carState.consumed = calculateFuelState(rate, carState.consumed, percent)
				if float32(size) >= carState.consumed {
					carState.fuel.Set(float32(size))
				} else {
					carState.consumed = float32(size)
					carState.fuel.Set(float32(size))
				}

				// TODO: make weight penalty configurable
				// c.maxSpeed[carId].Subtract = types.Percent(math.Round(float64(c.fuel[carId].Get() / 10)))
				// car.SetMaxSpeed().Update()

				/*				log.WithField("car", carId).
								WithField("fuel", c.fuel[carId].Get()).
								WithField("consumed", c.consumed[carId].Subtract).
								Debug("report car fuel level")
				*/
			}
		}
		return nil
	}
}

// calculateFuelState calculates the updated fuel state based on the burn rate, previous consumption, and trigger value.
// It returns the new fuel state.
func calculateFuelState(burnRate float32, consumed float32, triggerValue uint8) float32 {
	return ((float32(triggerValue) / 100) * burnRate) + consumed
}
