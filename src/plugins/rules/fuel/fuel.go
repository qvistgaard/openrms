package fuel

import (
	"openrms/implement"
	"openrms/plugins/rules/limbmode"
	"openrms/state"
)

const (
	defaultBurnRate = float32(100)
	defaultFuel     = float32(100)

	fuelState     = "fuel"
	burnRateState = "burnrate"
)

type Consumption struct {
}

func (c *Consumption) Notify(v *state.Value) {
	car, err := v.Owner().(state.Car)
	if !err {
		if v.Name() == state.CarEvent {
			event := v.Get().(implement.Event)
			if event.Ontrack {
				fs := car.State().Get(fuelState).Get().(float32)
				bs := car.State().Get(burnRateState).Get().(float32)

				cf := calculateFuelState(bs, fs, event.TriggerValue)
				car.State().Get(fuelState).Set(cf)

				if cf <= 0 {
					car.State().Get(limbmode.LimbMode).Set(true)
				}
			}
		}
	}
}

func (c *Consumption) InitializeCarState(car *state.Car) {
	f := car.State().Get(fuelState).Get()
	b := car.State().Get(burnRateState).Get()

	if f == nil {
		car.State().Get(fuelState).Set(defaultFuel)
	}
	if b == nil {
		car.State().Get(burnRateState).Set(defaultBurnRate)
	}
	car.State().Get(state.CarEvent).Subscribe(c)
}

func calculateFuelState(burnRate float32, fuel float32, triggerValue uint8) float32 {
	return fuel - ((float32(triggerValue) / 255) * burnRate)
}
