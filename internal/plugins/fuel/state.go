package fuel

import (
	"context"
	"github.com/qmuntal/stateless"
	"reflect"
)

const (
	stateCarOnTrack   = "car-on-track"
	stateCarDeslotted = "car-deslotted"

	triggerCarOnTrack      = "trigger-car-on-track"
	triggerCarDeslotted    = "trigger-car-deslotted"
	triggerUpdateFuelLevel = "trigger-update-fuel-level"
)

func machine(fuelUpdate func(ctx context.Context, args ...interface{}) error) *stateless.StateMachine {
	m := stateless.NewStateMachineWithMode(stateCarOnTrack, stateless.FiringImmediate)
	m.SetTriggerParameters(triggerUpdateFuelLevel, reflect.TypeOf(uint8(0)))
	m.Configure(stateCarOnTrack).
		InternalTransition(triggerUpdateFuelLevel, fuelUpdate).
		Permit(triggerCarDeslotted, stateCarDeslotted)

	m.Configure(stateCarDeslotted).
		Permit(triggerCarOnTrack, stateCarOnTrack)
	return m
}
