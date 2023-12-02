package pit

import (
	"context"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type MachineTrigger string
type MachineTriggerFunc func(trigger MachineTrigger) error

const (
	triggerCarEnteredPitLane       = MachineTrigger("CarEnteredPitLane")
	triggerCarExitedPitLane        = MachineTrigger("CarExitedPitLane")
	triggerCarStopped              = MachineTrigger("CarStopped")
	triggerCarMoving               = MachineTrigger("CarMoving")
	triggerCarPitStopConfirmed     = MachineTrigger("CarPitStopConfirmed")
	triggerCarPitStopAutoConfirmed = MachineTrigger("CarPitStopAutoConfirmed")
	triggerCarPitStopComplete      = MachineTrigger("CarPitStopComplete")
)

const (
	stateCarNotInPitLane    = "CarNotInPit"
	stateCarInPitLane       = "CarInPit"
	stateCarStopped         = "CarStopped"
	stateCarMoving          = "CarMoving"
	stateCarPitStopActive   = "CarPitStopActive"
	stateCarPitStopComplete = "CarPitStopComplete"
)

type CarPitState uint8

const (
	PitStateNotInPitLane CarPitState = iota
	PitStateEntered
	PitStateWaiting
	PitStateActive
	PitStateComplete
)

func alwaysIgnoreTrigger(context.Context, ...interface{}) bool {
	return true
}

func logPitStateChange(carId types.CarId, state string, logline string) {
	log.WithField("car", carId).WithField("state", state).Info(logline)
}

func logPitStateChangeAction(carId types.CarId, state string, logline string) stateless.ActionFunc {
	return func(ctx context.Context, args ...interface{}) error {
		logPitStateChange(carId, state, logline)
		return nil
	}
}

func machine(h Handler) *stateless.StateMachine {
	carId := h.Id()
	m := stateless.NewStateMachineWithMode(stateCarNotInPitLane, stateless.FiringImmediate)
	m.Configure(stateCarNotInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarNotInPitLane, "car exited pit lane")).
		Permit(triggerCarEnteredPitLane, stateCarInPitLane)
	// Ignore(triggerCarMoving, alwaysIgnoreTrigger).
	// Ignore(triggerCarStopped, alwaysIgnoreTrigger).
	// Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
	// Ignore(triggerCarExitedPitLane, alwaysIgnoreTrigger)

	m.Configure(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarInPitLane, "car entered pit lane")).
		// Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarMoving, stateCarMoving)
	// Permit(triggerCarStopped, stateCarStopped)

	m.Configure(stateCarMoving).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarMoving, "car is moving inside the pit lane")).
		// Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		// Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped)

	m.Configure(stateCarStopped).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarStopped, "car stopped inside the pit lane")).
		OnEntry(handleOnCarStop(m, h)).
		OnExit(handleOnCarStart(m, h)).
		//	Ignore(triggerCarStopped, alwaysIgnoreTrigger)
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarPitStopConfirmed, stateCarPitStopActive).
		// Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarPitStopAutoConfirmed, stateCarPitStopActive)

	m.Configure(stateCarPitStopActive).
		SubstateOf(stateCarStopped).
		Permit(triggerCarPitStopComplete, stateCarPitStopComplete).
		// Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		// Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		// Ignore(triggerCarPitStopAutoConfirmed, alwaysIgnoreTrigger).
		// Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		OnEntry(logPitStateChangeAction(carId, stateCarPitStopActive, "entering active pit state")).
		OnEntry(startPitStop(m, h))

	m.Configure(stateCarPitStopComplete).
		SubstateOf(stateCarPitStopActive).
		OnEntry(logPitStateChangeAction(carId, stateCarPitStopActive, "Pit stop complete")).
		OnEntry(handleOnOnComplete(h)).
		Permit(triggerCarMoving, stateCarMoving)
	// Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
	// Permit(triggerCarStopped, stateCarStopped)

	return m
}

func startPitStop(m *stateless.StateMachine, h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.Start(func(trigger MachineTrigger) error {
			return m.Fire(trigger)
		})
	}
}

func handleOnCarStop(machine *stateless.StateMachine, h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.OnCarStop(func(trigger MachineTrigger) error {
			return machine.Fire(trigger)
		})
	}
}

func handleOnCarStart(_ *stateless.StateMachine, h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.OnCarStart()
	}
}

func handleOnOnComplete(h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.OnComplete()
	}
}
