package fuel

import (
	"context"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state/rx/car"
	"github.com/qvistgaard/openrms/internal/state/rx/race"
	"github.com/qvistgaard/openrms/internal/state/rx/rules"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"reflect"
	"time"
)

const (
	stateCarOnTrack   = "carOnTrack"
	stateCarDeslotted = "carDeslotted"
)

const (
	triggerCarOnTrack      = "carOnTrack"
	triggerCarDeslotted    = "carDeslotted"
	triggerUpdateFuelLevel = "updateFuelLevels"
)

const (
	// Using Lemans fuel rules for normal petrol cars burn rate it 110kg/h
	// that leads to 0.03 kg pr second, petrol has a weight around 775 gr / liter
	// that means the burn it is about 0.023 l/per second scaling that by the random number
	// Gotten from a internet forum about scale models and wind tunnel testing (5.65) we get the
	// burn rate.
	// defaultBurnRate = types.LiterPerSecond(0.223) // / 5.65)
	defaultBurnRate = types.LiterPerSecond(40) // / 5.65)

	// LMP1 fuel tank size is 75 Liters
	defaultFuel     = types.Liter(75)
	defaultFlowRate = types.LiterPerSecond(2 * 5.65)

	CarFuel           = "car-fuel"
	CarConfigFlowRate = "car-config-flow-rate"
	CarConfigFuel     = "car-config-fuel"
	CarConfigBurnRate = "car-config-fuel-burn-rate"
)

type Consumption struct {
	fuel       map[types.Id]*reactive.Liter
	consumed   map[types.Id]*reactive.LiterSubtractModifier
	state      map[types.Id]*stateless.StateMachine
	config     *Config
	rules      rules.Rules
	raceStatus implement.RaceStatus
}

func (c *Consumption) Priority() int {
	return 10
}

func (c *Consumption) Name() string {
	return "fuel"
}

func (c *Consumption) HandlePitStop(car *car.Car, cancel <-chan bool) bool {
	log.WithField("car", car.Id()).
		WithField("fuel", c.fuel[car.Id()].Get()).
		Infof("fuel: refuelling started")
	for {
		select {
		case v := <-cancel:
			log.WithField("car", car.Id()).
				WithField("cancel", v).
				WithField("length", len(cancel)).
				Info("fuel: refuelling cancelled")
			return false
		case <-time.After(250 * time.Millisecond):
			used, full := calculateRefuellingValue(c.consumed[car.Id()].Subtract, defaultFlowRate/4)

			log.WithField("car", car.Id()).
				WithField("fuel-used", used).
				WithField("fuel", c.fuel[car.Id()].Get()).
				WithField("length", len(cancel)).
				Info("fuel: refuelling")

			c.consumed[car.Id()].Subtract = used
			c.fuel[car.Id()].Update()

			if full {
				log.WithField("car", car.Id()).
					WithField("fuel-used", used).
					WithField("fuel", c.fuel[car.Id()].Get()).
					Info("fuel: refuelling complete")
				return true
			}
		}
	}
}

func (c *Consumption) ConfigureCarState(car *car.Car) {
	a := reactive.Annotations{
		annotations.CarId: car.Id(),
	}
	c.consumed[car.Id()] = &reactive.LiterSubtractModifier{Subtract: 0}
	c.fuel[car.Id()] = reactive.NewLiter(50, a, reactive.Annotations{annotations.CarValueFieldName: fields.Fuel})
	c.fuel[car.Id()].Modifier(c.consumed[car.Id()], 1000)

	machine := stateless.NewStateMachineWithMode(stateCarOnTrack, stateless.FiringImmediate)
	machine.SetTriggerParameters(triggerUpdateFuelLevel, reflect.TypeOf(types.Percent(0)))
	machine.Configure(stateCarOnTrack).
		InternalTransition(triggerUpdateFuelLevel, func(ctx context.Context, args ...interface{}) error {
			percent := args[0].(types.Percent)
			if percent > 0 {
				liter := c.fuel[car.Id()].Get()
				if liter > 0 {
					c.consumed[car.Id()].Subtract = calculateFuelState(defaultBurnRate, c.consumed[car.Id()].Subtract, percent)
					c.fuel[car.Id()].Update()
					log.WithField("car", car.Id()).
						WithField("fuel", c.fuel[car.Id()].Get()).
						Trace("report car fuel level")
				}
			}
			return nil
		}).
		Permit(triggerCarDeslotted, stateCarDeslotted)

	machine.Configure(stateCarDeslotted).
		Permit(triggerCarOnTrack, stateCarOnTrack)

	c.state[car.Id()] = machine

	car.Deslotted().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			if i.(bool) {
				machine.Fire(triggerCarDeslotted)
			} else {
				machine.Fire(triggerCarOnTrack)
			}
		})
	})
	car.Controller().TriggerValue().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			machine.Fire(triggerUpdateFuelLevel, i.(types.Percent))
		})
	})

	c.fuel[car.Id()].RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			p := i.(types.Liter)
			if p <= 0 {
				mode := c.rules.CarRule("limb-mode").(*limbmode.LimbMode)
				if mode != nil {
					log.Info("Car ran out of fuel, enable limb-mode")
					mode.Enable(car)
				}
			}
		})
	})
}

func (c *Consumption) ConfigureRaceState(race *race.Race) {
	race.Status().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			c.raceStatus = i.(implement.RaceStatus)
			if c.raceStatus == implement.RaceStopped {
				for id, v := range c.consumed {
					v.Enabled = false
					c.fuel[id].Update()
				}
			}
			if c.raceStatus == implement.RaceRunning {
				for id, v := range c.consumed {
					v.Enabled = true
					v.Subtract = 0
					c.fuel[id].Update()
				}
			}
		})
	})
}

func (c *Consumption) InitializeRaceState(*race.Race, context.Context, reactive.ValuePostProcessor) {

}

func (c *Consumption) InitializeCarState(car *car.Car, ctx context.Context, postProcess reactive.ValuePostProcessor) {
	c.fuel[car.Id()].Init(ctx, postProcess)
	c.fuel[car.Id()].Update()
}

func calculateRefuellingValue(used types.Liter, flowRate types.LiterPerSecond) (types.Liter, bool) {
	liter := used - types.Liter(flowRate)
	if liter <= 0 {
		return 0, true
	} else {
		return liter, false
	}
}

func calculateFuelState(burnRate types.LiterPerSecond, fuel types.Liter, triggerValue types.Percent) types.Liter {
	used := (float64(triggerValue) / 100) * float64(burnRate)
	return types.Liter(used) + fuel
}
