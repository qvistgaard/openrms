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
	stateCarInPit     = machineState("car-in-pit")

	triggerCarOnTrack      = machineTrigger("trigger-car-on-track")
	triggerCarDeslotted    = machineTrigger("trigger-car-deslotted")
	triggerCarInPit        = machineTrigger("trigger-in-pit")
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
