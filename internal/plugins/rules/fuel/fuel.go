package fuel

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	"time"
)

type Liter float32
type LiterPerSecond float32

const (
	defaultBurnRate   = LiterPerSecond(100)
	defaultFuel       = Liter(90)
	defaultRefuelRate = LiterPerSecond(2)

	CarFuel           = "car-fuel"
	CarConfigFuel     = "car-config-fuel"
	CarConfigBurnRate = "car-config-fuel-burn-rate"
)

type Consumption struct {
}

func (c *Consumption) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if v.Name() == state.CarEventSequence && c.Get(state.CarOnTrack).(bool) {
			fs := c.Get(CarFuel).(Liter)
			bs := c.Get(CarConfigBurnRate).(LiterPerSecond)
			cf := calculateFuelState(bs, fs, v.Get().(uint))

			if cf <= 0 {
				c.Set(limbmode.CarLimbMode, true)
				c.Set(CarFuel, Liter(0))
			} else {
				c.Set(CarFuel, cf)

			}
		}
	}
}

func (c *Consumption) InitializeRaceState(race *state.Course) {

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
		if v >= m {
			v = m
		}
		car.Set(CarFuel, v)
		return
	}
}

func (c *Consumption) Priority() uint8 {
	return 1
}

func calculateFuelState(burnRate LiterPerSecond, fuel Liter, triggerValue uint) Liter {
	return Liter(float32(fuel) - ((float32(triggerValue) / 255) * float32(burnRate)))
}
