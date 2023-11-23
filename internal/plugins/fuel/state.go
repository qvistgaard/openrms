package fuel

import (
	"context"
	"github.com/qmuntal/stateless"
	"reflect"
)

type machineState string
type machineTrigger string

const (
	stateCarOnTrack   = machineState("car-on-track")
	stateCarDeslotted = machineState("car-deslotted")

	triggerCarOnTrack      = machineTrigger("trigger-car-on-track")
	triggerCarDeslotted    = machineTrigger("trigger-car-deslotted")
	triggerUpdateFuelLevel = machineTrigger("trigger-update-fuel-level")
)

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
				// car.MaxSpeed().Update()

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
