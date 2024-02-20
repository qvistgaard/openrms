package pit

import (
	"context"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type MachineTrigger string
type MachineState string
type StartPitStop func() error
type CancelPitStop func() error
type CompletePitStop func() error

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
	stateCarNotInPitLane    = MachineState("CarNotInPit")
	stateCarInPitLane       = MachineState("CarInPit")
	stateCarStopped         = MachineState("CarStopped")
	stateCarMoving          = MachineState("CarMoving")
	stateCarPitStopActive   = MachineState("CarPitStopActive")
	stateCarPitStopComplete = MachineState("CarPitStopComplete")
)

func logPitStateChange(carId types.CarId, state MachineState, logline string) {
	log.WithField("car", carId).WithField("state", state).Info(logline)
}

func logPitStateChangeAction(carId types.CarId, state MachineState, logline string) stateless.ActionFunc {
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

	m.Configure(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarInPitLane, "car entered pit lane")).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane)

	m.Configure(stateCarMoving).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarMoving, "car is moving inside the pit lane")).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped).
		Ignore(triggerCarMoving)

	m.Configure(stateCarStopped).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(carId, stateCarStopped, "car stopped inside the pit lane")).
		OnEntry(handleOnCarStop(m, h)).
		OnExit(handleOnCarStart(m, h)).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarPitStopConfirmed, stateCarPitStopActive).
		Permit(triggerCarPitStopAutoConfirmed, stateCarPitStopActive).
		Ignore(triggerCarStopped)

	m.Configure(stateCarPitStopActive).
		SubstateOf(stateCarStopped).
		Permit(triggerCarPitStopComplete, stateCarPitStopComplete).
		OnEntry(logPitStateChangeAction(carId, stateCarPitStopActive, "entering active pit state")).
		OnEntry(startPitStop(m, h))

	m.Configure(stateCarPitStopComplete).
		OnEntry(logPitStateChangeAction(carId, stateCarPitStopActive, "Pit stop complete")).
		OnEntry(handleOnOnComplete(h)).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane)

	return m
}

func startPitStop(m *stateless.StateMachine, h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.Start(func() error {
			return m.Fire(triggerCarPitStopComplete)
		}, func() error {
			return m.Fire(triggerCarPitStopComplete)
		})
	}
}

func handleOnCarStop(m *stateless.StateMachine, h Handler) func(ctx context.Context, args ...any) error {
	return func(ctx context.Context, args ...any) error {
		return h.OnCarStop(func() error {
			return m.Fire(triggerCarPitStopConfirmed)
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
