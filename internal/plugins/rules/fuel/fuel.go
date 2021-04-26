package fuel

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	"time"
)

type Liter float32
type LiterPerSecond float32

const (
	defaultBurnRate   = LiterPerSecond(0.1)
	defaultFuel       = Liter(90)
	defaultRefuelRate = LiterPerSecond(2)

	CarFuel           = "car-fuel"
	CarConfigFuel     = "car-config-fuel"
	CarConfigBurnRate = "car-config-fuel-burn-rate"
)

type Consumption struct {
	course *state.Course
}

func (c *Consumption) Notify(v *state.Value) {
	if car, ok := v.Owner().(*state.Car); ok {
		if v.Name() == state.CarEventSequence && car.Get(state.CarOnTrack).(bool) {
			if rs, ok := c.course.Get(state.RaceStatus).(uint8); !ok || rs != state.RaceStatusPaused {
				fs := car.Get(CarFuel).(Liter)
				bs := car.Get(CarConfigBurnRate).(LiterPerSecond)
				tv := car.Get(state.ControllerTriggerValue).(state.TriggerValue)
				cf := calculateFuelState(bs, fs, tv)

				if cf <= 0 {
					car.Set(limbmode.CarLimbMode, true)
					car.Set(CarFuel, Liter(0))
				} else {
					car.Set(CarFuel, cf)
				}
			}
		}
	}
}

func (c *Consumption) InitializeCourseState(race *state.Course) {
	c.course = race
}

func (c *Consumption) InitializeCarState(car *state.Car) {
	f := car.Get(CarFuel)
	cf := car.Get(CarConfigFuel)
	cb := car.Get(CarConfigBurnRate)

	if cf == nil {
		car.Set(CarConfigFuel, defaultFuel)
	}
	if f == nil {
		car.Set(CarFuel, car.Get(CarConfigFuel))
	}
	if cb == nil {
		car.Set(CarConfigBurnRate, defaultBurnRate)
	}
	car.Subscribe(state.CarEventSequence, c)
}

func (c *Consumption) HandlePitStop(car *state.Car, cancel chan bool) {
	select {
	case <-cancel:
		return
	case <-time.After(500 * time.Millisecond):
		f := car.Get(CarFuel).(Liter)
		v := f + Liter(defaultRefuelRate/2)
		m := car.Get(CarConfigFuel).(Liter)
		d := false
		if v >= m {
			v = m
			d = true

		}
		car.Set(CarFuel, v)
		if d {
			return
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
