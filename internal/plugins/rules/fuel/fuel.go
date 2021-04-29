package fuel

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"time"
)

type Liter float32
type LiterPerSecond float32

const (
	// Using Lemans fuel rules for normal petrol cars burn rate it 110kg/h
	// that leads to 0.03 kg pr second, petrol has a weight around 775 gr / liter
	// that means the burn it is about 0.023 l/per second scaling that by the random number
	// Gotten from a internet forum about scale models and wind tunnel testing (5.65) we get the
	// burn rate.
	defaultBurnRate = LiterPerSecond(0.023 / 2) // / 5.65)

	// LMP1 fuel tank size is 75 Liters
	defaultFuel     = Liter(75)
	defaultFlowRate = LiterPerSecond(2 * 5.65)

	CarFuel           = "car-fuel"
	CarConfigFlowRate = "car-config-flow-rate"
	CarConfigFuel     = "car-config-fuel"
	CarConfigBurnRate = "car-config-fuel-burn-rate"
)

type Consumption struct {
	course *state.Course
	config *Config
}

func (c *Consumption) Notify(v *state.Value) {
	log.Infof("state: %+v", c.course.Get(state.RaceStatus))
	if c.course.Get(state.RaceStatus) != state.RaceStatusStopped {
		if car, ok := v.Owner().(*state.Car); ok {
			if v.Name() == state.CarEventSequence && car.Get(state.CarOnTrack).(bool) {
				if rs, ok := c.course.Get(state.RaceStatus).(uint8); !ok || rs != state.RaceStatusPaused {
					fs := car.Get(CarFuel).(Liter)
					bs := car.Get(CarConfigBurnRate).(LiterPerSecond)
					tv := car.Get(state.ControllerTriggerValue).(state.TriggerValue)
					cf := calculateFuelState(bs, fs, tv)

					if cf <= 0 {
						log.WithField("car", car.Id()).Info("car has run out of fuel.")
						car.Set(limbmode.CarLimbMode, true)
						car.Set(CarFuel, Liter(0))
					} else {
						car.Set(CarFuel, cf)
					}
				}
			}
		}
	}
}

func (c *Consumption) InitializeCourseState(course *state.Course) {
	c.course = course
}

func (c *Consumption) InitializeCarState(car *state.Car) {
	cc := &Config{}
	err := car.Settings(cc)
	if err != nil {
		log.Fatal(err)
	}

	var fuel Liter
	if cc.Fuel != nil {
		fuel = *cc.Fuel
	} else if c.config.Fuel != nil {
		fuel = *c.config.Fuel
	} else {
		fuel = defaultFuel
	}
	car.Set(CarConfigFuel, fuel)

	var startingFuel Liter
	if cc.StartingFuel != nil {
		startingFuel = *cc.StartingFuel
	} else if c.config.StartingFuel != nil {
		startingFuel = *c.config.StartingFuel
	} else {
		startingFuel = fuel
	}
	car.Set(CarFuel, startingFuel)

	var burnRate LiterPerSecond
	if cc.BurnRate != nil {
		burnRate = *cc.BurnRate
	} else if c.config.BurnRate != nil {
		burnRate = *c.config.BurnRate
	} else {
		burnRate = defaultBurnRate
	}
	car.Set(CarConfigBurnRate, burnRate)

	var flowRate LiterPerSecond
	if c.config.FlowRate != nil {
		flowRate = *c.config.FlowRate
	} else {
		flowRate = defaultFlowRate
	}
	car.Set(CarConfigFlowRate, flowRate)

	car.Subscribe(state.CarEventSequence, c)
}

func (c *Consumption) HandlePitStop(car *state.Car, cancel chan bool) {
	log.WithField("car", car.Id()).Infof("fuel: refuelling started")
	for {
		select {
		case <-cancel:
			log.WithField("car", car.Id()).Infof("fuel: refuelling cancelled")
			return
		case <-time.After(250 * time.Millisecond):
			f := car.Get(CarFuel).(Liter)
			v := f + Liter(car.Get(CarConfigFlowRate).(LiterPerSecond)/4)
			m := car.Get(CarConfigFuel).(Liter)
			d := false
			if v >= m {
				v = m
				d = true
			}
			car.Set(CarFuel, v)
			if d {
				log.WithField("car", car.Id()).Infof("fuel: refuelling complete")
				return
			}
		}
	}
}

func (c *Consumption) Priority() uint8 {
	return 1
}

func calculateFuelState(burnRate LiterPerSecond, fuel Liter, triggerValue state.TriggerValue) Liter {
	used := float32(triggerValue) * float32(burnRate)
	remaining := float32(fuel) - used

	if remaining > 0 {
		return Liter(remaining)
	} else {
		return Liter(0)
	}
}
