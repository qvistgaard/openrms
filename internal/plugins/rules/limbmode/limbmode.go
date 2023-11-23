package limbmode

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
)

const CarLimbMode = "limb-mode"
const CarLimbModeMaxSpeed = "limb-mode-max-speedModifier"

type Settings struct {
	LimbMode struct {
		MaxSpeed *types.Percent `mapstructure:"max-speedModifier,omitempty"`
	} `mapstructure:"limb-mode"`
}

type LimbMode struct {
	defaults *LimbModeConfig
	config   map[types.Id]*LimbModeConfig
	state    map[types.Id]observable.Observable[bool]
	// speedModifier map[types.Id]*reactive.PercentAbsoluteModifier
}

func (l *LimbMode) ConfigureCarState(car *car.Car) {
	// var carConfig *LimbModeConfig
	// var ok bool
	/*	if carConfig, ok = l.config[car.Id()]; !ok {
		carConfig = l.defaults
	}*/
	a := []observable.Annotation{
		{annotations.CarId, car.Id().String()},
	}

	l.state[car.Id()] = observable.Create(false, append(a, observable.Annotation{annotations.CarValueFieldName, "limb-mode"})...)
	// l.speedModifier[car.Id()] = &reactive.PercentAbsoluteModifier{Absolute: *carConfig.MaxSpeed, Enabled: false, Condition: reactive.IfGreaterThen}
	// car.MaxSpeed().Modifier(l.speedModifier[car.Id()], 1)

	/*	l.state[car.Id()].RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			l.speedModifier[car.Id()].Enabled = i.(bool)
			//	car.MaxSpeed().Update()
		})
	})*/
}

func (l *LimbMode) ConfigureRaceState(raceState *race.Race) {
	/*	raceState.Status().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			s := i.(race.RaceStatus)
			if s == race.RaceStopped {
				for id, v := range l.speedModifier {
					v.Enabled = false
					l.state[id].Update()
				}
			}
		})
	})*/
}

func (l *LimbMode) HandlePitStop(car *car.Car, cancel <-chan bool) bool {
	/*
		if l.speedModifier[car.Id()].Enabled {
			log.WithField("car", car.Id()).
				Info("limb-mode penalty started")
			select {
			case <-time.After(5000 * time.Millisecond):
				log.WithField("car", car.Id()).
					Info("limb-mode penalty complete")
				l.Disable(car)
			}
		}

	*/
	return true
}

func (l *LimbMode) InitializeRaceState(race *race.Race, ctx context.Context) {

}

func (l *LimbMode) InitializeCarState(car *car.Car, ctx context.Context) {
	// l.state[car.Id()].Init(ctx)
}

func (l *LimbMode) Get(car *car.Car) observable.Observable[bool] {
	return l.state[car.Id()]
}

func (l *LimbMode) Enable(car *car.Car) {
	l.state[car.Id()].Set(true)
	// l.state[car.Id()].Update()
}

func (l *LimbMode) Disable(car *car.Car) {
	l.state[car.Id()].Set(false)
	// l.state[car.Id()].Update()
}

func (l *LimbMode) Priority() int {
	return 1
}

func (l *LimbMode) Name() string {
	return "limb-mode"
}
