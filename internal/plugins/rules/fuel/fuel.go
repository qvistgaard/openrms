package fuel

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
)

const (
	defaultBurnRate = float32(100)
	defaultFuel     = float32(100)

	CarFuel           = "car-fuel"
	CarConfigFuel     = "car-config-fuel"
	CarConfigBurnRate = "car-config-fuel-burn-rate"
)

type Consumption struct {
}

func (c *Consumption) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		if v.Name() == state.CarEventSequence && c.Get(state.CarOnTrack).(bool) {
			fs := c.Get(CarFuel).(float32)
			bs := c.Get(CarConfigBurnRate).(float32)
			cf := calculateFuelState(bs, fs, v.Get().(uint8))

			c.Set(CarFuel, cf)
			if cf <= 0 {
				c.Set(limbmode.CarLimbMode, true)
			}
		}
	}
}

func (c *Consumption) InitializeRaceState(race *state.Race) {

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

func calculateFuelState(burnRate float32, fuel float32, triggerValue uint8) float32 {
	return fuel - ((float32(triggerValue) / 255) * burnRate)
}
