package pit

import (
	"context"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/state/rx/car"
	"github.com/qvistgaard/openrms/internal/state/rx/rules"
	"github.com/qvistgaard/openrms/internal/types"
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

type Rule struct {
	rules            rules.Rules
	speed            map[types.Id]*reactive.PercentSubtractModifier
	stops            map[types.Id]chan bool
	carState         map[types.Id]*stateless.StateMachine
	handlerCompleted map[types.Id]bool

	// course   *state.Course
}

func CreatePitRule(rules rules.Rules) *Rule {
	p := new(Rule)
	p.rules = rules
	p.stops = make(map[types.Id]chan bool)
	p.speed = make(map[types.Id]*reactive.PercentSubtractModifier)
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
	stateMachine := p.newState(c)
	p.carState[c.Id()] = stateMachine
	p.handlerCompleted[c.Id()] = false
	p.speed[c.Id()] = &reactive.PercentSubtractModifier{Subtract: 100}
	p.stops[c.Id()] = make(chan bool, 10)

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
			if i.(bool) && !p.handlerCompleted[c.Id()] {
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

func (p *Rule) InitializeCarState(*car.Car, context.Context, reactive.ValuePostProcessor) {

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
			go pitStopActivationHandler(car, machine, cancel)
		}
		return nil
	}
}

func pitStopActivationHandler(car *car.Car, machine *stateless.StateMachine, cancel chan bool) {
	log.WithField("car", car.Id()).Debug("waiting for pit stop confirmation")
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
			p.handlerCompleted[car.Id()] = false
			return nil
		}).
		Permit(triggerCarEnteredPitLane, stateCarInPitLane).
		Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarExitedPitLane, alwaysIgnoreTrigger)

	machine.Configure(stateCarInPitLane).
		OnEntry(logPitStateChangeAction(car, stateCarInPitLane, "car entered pit lane")).
		OnExit(cancelPitStopAutoConfirmation(cancelWait)).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarStopped, stateCarStopped)

	machine.Configure(stateCarMoving).
		OnEntry(logPitStateChangeAction(car, stateCarStopped, "car is moving inside the pit lane")).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarMoving, alwaysIgnoreTrigger).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped)

	machine.Configure(stateCarStopped).
		OnEntry(logPitStateChangeAction(car, stateCarStopped, "car is stopped inside the pit lane")).
		OnEntry(p.pitStopActivationHandlerAction(car, machine, cancelWait)).
		OnExit(cancelPitStopAutoConfirmation(cancelWait)).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarPitStopConfirmed, stateCarPitStopActive, func(ctx context.Context, args ...interface{}) bool {
			return !p.handlerCompleted[car.Id()]
		}).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarPitStopAutoConfirmed, stateCarPitStopActive, func(ctx context.Context, args ...interface{}) bool {
			return !p.handlerCompleted[car.Id()]
		})

	machine.Configure(stateCarPitStopActive).
		OnEntry(logPitStateChangeAction(car, stateCarPitStopActive, "entering active pit state")).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			p.speed[car.Id()].Enabled = true
			car.MaxSpeed().Update()
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
			car.MaxSpeed().Update()
			log.Info("re-enable car")
			return nil
		}).
		Permit(triggerCarPitStopComplete, stateCarPitStopComplete).
		Ignore(triggerCarPitStopConfirmed, alwaysIgnoreTrigger).
		Ignore(triggerCarStopped, alwaysIgnoreTrigger).
		Ignore(triggerCarPitStopAutoConfirmed, alwaysIgnoreTrigger)

	machine.Configure(stateCarPitStopComplete).
		OnEntry(logPitStateChangeAction(car, stateCarPitStopComplete, "pit stop complete")).
		Permit(triggerCarMoving, stateCarMoving).
		Permit(triggerCarExitedPitLane, stateCarNotInPitLane).
		Permit(triggerCarStopped, stateCarStopped)

	return machine
}
