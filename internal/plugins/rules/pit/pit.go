package pit

import (
	"context"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/rules"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	triggerCarEnteredPitLane       = "CarEnteredPitLane"
	triggerCarExitedPitLane        = "CarExitedPitLane"
	triggerCarStopped              = "CarStopped"
	triggerCarMoving               = "CarMoving"
	triggerCarPitStopConfirmed     = "CarPitStopConfirmed"
	triggerCarPitStopAutoConfirmed = "CarPitStopAutoConfirmed"
	triggerCarPitStopComplete      = "CarPitStopComplete"
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

type CarPitStateValue struct {
	reactive.Value
}

func NewCarPitStateValue(initial CarPitState, annotations ...reactive.Annotations) *CarPitStateValue {
	return &CarPitStateValue{reactive.NewDistinctValue(initial, annotations...)}
}

func (p *CarPitStateValue) Set(value CarPitState) {
	p.Value.Set(value)
}

type Rule struct {
	rules            rules.Rules
	speed            map[types.Id]*reactive.PercentSubtractModifier
	carState         map[types.Id]*stateless.StateMachine
	carPitState      map[types.Id]*CarPitStateValue
	handlerCompleted map[types.Id]bool

	// course   *state.Course
}

func CreatePitRule(rules rules.Rules) *Rule {
	p := new(Rule)
	p.rules = rules
	p.speed = make(map[types.Id]*reactive.PercentSubtractModifier)
	p.carPitState = make(map[types.Id]*CarPitStateValue)
	p.handlerCompleted = make(map[types.Id]bool)
	p.carState = make(map[types.Id]*stateless.StateMachine)
	return p
}

func (p *Rule) Priority() int {
	return 9
}

func (p *Rule) Name() string {
	return "pit"
}

func (p *Rule) ConfigureCarState(c *car.Car) {
	carId := c.Id()
	a := reactive.Annotations{
		annotations.CarId: carId,
	}

	stateMachine := p.newState(c)
	p.carState[carId] = stateMachine
	p.handlerCompleted[carId] = false
	p.speed[carId] = &reactive.PercentSubtractModifier{Subtract: 100}
	p.carPitState[carId] = NewCarPitStateValue(PitStateNotInPitLane, a, reactive.Annotations{annotations.CarValueFieldName: fields.PitState})

	c.PitLaneMaxSpeed().Modifier(p.speed[carId], 1000)

	c.Pit().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			var err error
			if b, ok := i.(bool); ok && !b {
				err = stateMachine.Fire(triggerCarExitedPitLane)
			} else {
				err = stateMachine.Fire(triggerCarEnteredPitLane)
			}
			if err != nil {
				log.Error(err)
			}
		})
	})
	c.Controller().ButtonTrackCall().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			if i.(bool) && !p.handlerCompleted[carId] {
				err := stateMachine.Fire(triggerCarPitStopConfirmed)
				if err != nil {
					log.Error(err)
				}
			}
		})
	})
	c.Controller().TriggerValue().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			var err error
			triggerValue := i.(types.Percent)
			if triggerValue == 0 {
				err = stateMachine.Fire(triggerCarStopped)
			} else {
				err = stateMachine.Fire(triggerCarMoving)
			}
			if err != nil {
				log.Error(err)
			}
		})
	})
}

func (p *Rule) InitializeCarState(car *car.Car, ctx context.Context, postProcess reactive.ValuePostProcessor) {
	p.carPitState[car.Id()].Init(ctx, postProcess)
}

func alwaysIgnoreTrigger(context.Context, ...interface{}) bool {
	return true
}

func logPitStateChange(car *car.Car, state string, logline string) {
	log.WithField("car", car.Id()).WithField("state", state).Info(logline)
}

func logPitStateChangeAction(car *car.Car, state string, logline string) stateless.ActionFunc {
	return func(ctx context.Context, args ...interface{}) error {
		logPitStateChange(car, state, logline)
		return nil
	}
}

func cancelPitStopAutoConfirmation(cancel chan bool) stateless.ActionFunc {
	return func(ctx context.Context, args ...interface{}) error {
		if len(cancel) == 0 {
			cancel <- true
		}
		return nil
	}
}

func (p *Rule) pitStopActivationHandlerAction(car *car.Car, machine *stateless.StateMachine, cancel chan bool) stateless.ActionFunc {
	return func(ctx context.Context, args ...interface{}) error {
		if !p.handlerCompleted[car.Id()] {
			p.carPitState[car.Id()].Set(PitStateWaiting)
			p.carPitState[car.Id()].Update()
			go pitStopActivationHandler(car, machine, cancel)
		}
		return nil
	}
}

func pitStopActivationHandler(car *car.Car, machine *stateless.StateMachine, cancel chan bool) {
	log.WithField("car", car.Id()).Debug("waiting for automatic pit stop confirmation")
	select {
	case <-time.After(5 * time.Second):
		log.WithField("car", car.Id()).Info("pit stop automatically confirmed")
		err := machine.Fire(triggerCarPitStopAutoConfirmed)
		if err != nil {
			log.Error(err)
		}
	case <-cancel:
		log.WithField("car", car.Id()).Debug("pit stop wait cancelled")
	}
}

func (p *Rule) newState(car *car.Car) *stateless.StateMachine {
	cancelWait := make(chan bool, 1)
	machine := stateless.NewStateMachineWithMode(stateCarNotInPitLane, stateless.FiringImmediate)
	machine.Configure(stateCarNotInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarNotInPitLane, "car exited pit lane")).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			for len(cancelWait) > 0 {
				<-cancelWait
			}

			p.handlerCompleted[car.Id()] = false
			p.carPitState[car.Id()].Set(PitStateNotInPitLane)
			p.carPitState[car.Id()].Update()
			return nil
		}).
		Permit(triggerCarEnteredPitLane, stateCarInPitLane).
		Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarExitedPitLane, alwaysIgnoreTrigger)

	machine.Configure(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarInPitLane, "car entered pit lane")).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			p.carPitState[car.Id()].Set(PitStateEntered)
			p.carPitState[car.Id()].Update()
			return nil
		}).
		OnExit(cancelPitStopAutoConfirmation(cancelWait)).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarStopped, stateCarStopped)

	machine.Configure(stateCarMoving).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarStopped, "car is moving inside the pit lane")).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped)

	machine.Configure(stateCarStopped).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarStopped, "car is stopped inside the pit lane")).
		OnEntry(p.pitStopActivationHandlerAction(car, machine, cancelWait)).
		OnExit(cancelPitStopAutoConfirmation(cancelWait)).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarPitStopConfirmed, stateCarPitStopActive, func(ctx context.Context, args ...interface{}) bool {
			log.Info("Manually triggered")
			return !p.handlerCompleted[car.Id()]
		}).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarPitStopAutoConfirmed, stateCarPitStopActive, func(ctx context.Context, args ...interface{}) bool {
			log.Info("Auto triggered")
			return !p.handlerCompleted[car.Id()]
		})

	machine.Configure(stateCarPitStopActive).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarPitStopActive, "entering active pit state")).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			p.carPitState[car.Id()].Set(PitStateActive)
			p.carPitState[car.Id()].Update()

			// TODO: Figure out why car can drive when pit is active

			p.speed[car.Id()].Enabled = true
			car.PitLaneMaxSpeed().Update()
			go func() {
				for _, r := range p.rules.PitRules() {
					if !r.HandlePitStop(car, make(chan bool)) {
						break
					}
				}
				log.Info("pit stop complete")
				err := machine.Fire(triggerCarPitStopComplete)
				if err != nil {
					log.Error(err)
				}
			}()
			return nil
		}).
		OnExit(func(ctx context.Context, args ...interface{}) error {
			p.handlerCompleted[car.Id()] = true
			p.speed[car.Id()].Enabled = false
			car.PitLaneMaxSpeed().Update()
			log.Info("re-enable car")
			return nil
		}).
		Permit(triggerCarPitStopComplete, stateCarPitStopComplete).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Ignore(triggerCarPitStopAutoConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarMoving, alwaysIgnoreTrigger)

	machine.Configure(stateCarPitStopComplete).
		SubstateOf(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarPitStopComplete, "pit stop complete")).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			p.carPitState[car.Id()].Set(PitStateComplete)
			p.carPitState[car.Id()].Update()
			return nil
		}).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped)

	return machine
}
