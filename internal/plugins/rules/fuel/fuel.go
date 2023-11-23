package fuel

import (
	"context"
	"github.com/divideandconquer/go-merge/merge"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/rules"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
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

type fuelConsumed struct {
	enabled  bool
	consumed float32
}

type Consumption struct {
	fuel     map[types.Id]observable.Observable[float32]
	consumed map[types.Id]fuelConsumed
	// maxSpeed   map[types.Id]*reactive.PercentSubtractModifier
	state      map[types.Id]*stateless.StateMachine
	config     *Config
	fuelConfig map[types.Id]*FuelConfig
	rules      rules.Rules
	raceStatus race.RaceStatus
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
			// used, full := calculateRefuellingValue(c.consumed[car.Id()].Subtract, c.fuelConfig[car.Id()].FlowRate/4)
			used := 1
			full := true

			log.WithField("car", car.Id()).
				WithField("fuel-used", used).
				WithField("fuel", c.fuel[car.Id()].Get()).
				WithField("length", len(cancel)).
				Trace("fuel: refuelling")

			// c.consumed[car.Id()].Subtract = used
			c.fuel[car.Id()].Publish()

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
	carId := car.Id()
	for _, v := range c.config.Car.Cars {
		if *v.Id == carId {
			c.fuelConfig[carId] = merge.Merge(c.config.Car.Defaults, v).(*CarSettings).FuelConfig
		}
	}
	a := []observable.Annotation{
		{annotations.CarId, carId.String()},
	}

	// c.consumed[carId] = &reactive.LiterSubtractModifier{Subtract: c.fuelConfig[carId].TankSize - c.fuelConfig[carId].StartingFuel, Enabled: true}
	c.fuel[carId] = observable.Create(float32(c.fuelConfig[carId].TankSize), append(a, observable.Annotation{annotations.CarValueFieldName, fields.Fuel})...)
	// c.fuel[carId] = valueFactory.NewLiter(c.fuelConfig[carId].TankSize, a, reactive.Annotations{annotations.CarValueFieldName: fields.Fuel})
	c.fuel[carId].Modifier(func(f float32) (float32, bool) {
		return f - c.consumed[carId].consumed, c.consumed[carId].enabled
	}, 1000)

	// c.fuel[carId].Modifier(c.consumed[carId], 1000)

	// c.maxSpeed[carId] = &reactive.PercentSubtractModifier{Subtract: 0, Enabled: true}
	// car.MaxSpeed().Modifier(c.maxSpeed[carId], 0)

	machine := stateless.NewStateMachineWithMode(stateCarOnTrack, stateless.FiringImmediate)
	machine.SetTriggerParameters(triggerUpdateFuelLevel, reflect.TypeOf(types.Percent(0)))
	machine.Configure(stateCarOnTrack).
		InternalTransition(triggerUpdateFuelLevel, c.handleUpdateFuelLevel(car, carId)).
		Permit(triggerCarDeslotted, stateCarDeslotted)

	machine.Configure(stateCarDeslotted).
		Permit(triggerCarOnTrack, stateCarOnTrack)

	c.state[carId] = machine

	/*	car.Deslotted().RegisterObserver(func(observable rxgo.Observable) {
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

		c.fuel[carId].RegisterObserver(func(observable rxgo.Observable) {
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
		})*/
}

func (c *Consumption) handleUpdateFuelLevel(car *car.Car, carId types.Id) func(ctx context.Context, args ...interface{}) error {
	return func(ctx context.Context, args ...interface{}) error {
		// trigger percentage
		percent := args[0].(types.Percent)
		if percent > 0 {
			liter := c.fuel[carId].Get()
			if liter > 0 {
				/*
					substract := calculateFuelState(c.fuelConfig[carId].BurnRate, c.consumed[carId].Subtract, percent)
					if c.fuelConfig[carId].TankSize >= substract {
						c.consumed[carId].Subtract = substract
						c.fuel[carId].Publish()
					}
				*/
				// TODO: make weight penalty configurable
				// c.maxSpeed[carId].Subtract = types.Percent(math.Round(float64(c.fuel[carId].Get() / 10)))
				// car.MaxSpeed().Update()

				/*
					log.WithField("car", carId).
						WithField("fuel", c.fuel[carId].Get()).
						WithField("consumed", c.consumed[carId].Subtract).
						Debug("report car fuel level")

				*/
			}
		}
		return nil
	}
}

func (c *Consumption) ConfigureRaceState(raceState *race.Race) {
	/*	raceState.Status().RegisterObserver(func(status race.RaceStatus) {
		c.raceStatus = status
		if c.raceStatus == race.RaceStopped {
			for id, v := range c.consumed {
				v.Enabled = false
				c.fuel[id].Update()
			}
		}
		if c.raceStatus == race.RaceRunning {
			for id, v := range c.consumed {
				v.Enabled = true
				v.Subtract = 0
				c.fuel[id].Update()
			}
		}
	})*/
}

func (c *Consumption) InitializeRaceState(*race.Race, context.Context) {

}

func (c *Consumption) InitializeCarState(car *car.Car, ctx context.Context) {
	// c.fuel[car.Id()].Init(ctx)
	// c.fuel[car.Id()].Update()
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
