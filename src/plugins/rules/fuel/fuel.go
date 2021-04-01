package fuel

import (
	"openrms/ipc"
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

/*func (c *Consumption) Handle(race connector.Connector, telemetry queue.Queue, car *state.Car, event ipc.Event) {
	if event.Ontrack {
		fs := car.Get(fuelState).Get().(float32)
		bs := car.Get(burnRateState).Get().(float32)

		cf := calculateFuelState(bs, fs, event.TriggerValue)
		car.Get(fuelState).Set(cf)

		if cf <= 0 {
			car.Get(limbmode.LimbMode).Set(true)
		}

		telemetry.Enqueue(telemetry2.Telemetry{Id: car.Id(), Name: fuelState, Value: cf, Time: time.Now()})
	}
}*/

func (c *Consumption) Notify(v *state.Value) {
	car, err := v.Owner().(state.Car)
	if !err {
		if v.Name() == state.RaceEvent {
			event := v.Get().(ipc.Event)
			if event.Ontrack {
				fs := car.Get(fuelState).Get().(float32)
				bs := car.Get(burnRateState).Get().(float32)

				cf := calculateFuelState(bs, fs, event.TriggerValue)
				car.Get(fuelState).Set(cf)

				if cf <= 0 {
					car.Get(limbmode.LimbMode).Set(true)
				}
			}
		}
	}
}

func (c *Consumption) InitializeCarState(car *state.Car) {
	f := car.Get(fuelState).Get()
	b := car.Get(burnRateState).Get()

	if f == nil {
		car.Get(fuelState).Set(defaultFuel)
	}
	if b == nil {
		car.Get(burnRateState).Set(burnRateState)
	}
	car.Get(state.RaceEvent).Subscribe(c)
}

func calculateFuelState(burnRate float32, fuel float32, triggerValue uint8) float32 {
	return fuel - ((float32(triggerValue) / 255) * burnRate)
}
