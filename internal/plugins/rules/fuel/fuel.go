package fuel

import (
	"context"
	"github.com/divideandconquer/go-merge/merge"
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

type Consumption struct {
	fuel       map[types.Id]*reactive.Liter
	consumed   map[types.Id]*reactive.LiterSubtractModifier
	state      map[types.Id]*stateless.StateMachine
	config     *Config
	fuelConfig map[types.Id]*FuelConfig
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
			used, full := calculateRefuellingValue(c.consumed[car.Id()].Subtract, c.fuelConfig[car.Id()].FlowRate/4)

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
	for _, v := range c.config.Car.Cars {
		if *v.Id == car.Id() {
			c.fuelConfig[car.Id()] = merge.Merge(c.config.Car.Defaults, v).(*CarSettings).FuelConfig
		}
	}
	a := reactive.Annotations{
		annotations.CarId: car.Id(),
	}
	c.consumed[car.Id()] = &reactive.LiterSubtractModifier{Subtract: c.fuelConfig[car.Id()].TankSize - c.fuelConfig[car.Id()].StartingFuel}
	c.fuel[car.Id()] = reactive.NewLiter(c.fuelConfig[car.Id()].TankSize, a, reactive.Annotations{annotations.CarValueFieldName: fields.Fuel})
	c.fuel[car.Id()].Modifier(c.consumed[car.Id()], 1000)

	machine := stateless.NewStateMachineWithMode(stateCarOnTrack, stateless.FiringImmediate)
	machine.SetTriggerParameters(triggerUpdateFuelLevel, reflect.TypeOf(types.Percent(0)))
	machine.Configure(stateCarOnTrack).
		InternalTransition(triggerUpdateFuelLevel, func(ctx context.Context, args ...interface{}) error {
			percent := args[0].(types.Percent)
			if percent > 0 {
				liter := c.fuel[car.Id()].Get()
				if liter > 0 {
					c.consumed[car.Id()].Subtract = calculateFuelState(c.fuelConfig[car.Id()].BurnRate, c.consumed[car.Id()].Subtract, percent)
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
