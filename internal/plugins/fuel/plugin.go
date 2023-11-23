package fuel

import (
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
)

type Plugin struct {
	fuel      map[types.Id]observable.Observable[float32]
	config    Config
	carConfig map[types.Id]CarSettings
	fuelState map[types.Id]fuelState
	status    race.RaceStatus
	limbMode  *limbmode.Plugin
}

type fuelState struct {
	enabled  bool
	consumed float32
}

func New(config Config, limbMode *limbmode.Plugin) (*Plugin, error) {
	return &Plugin{
		fuel:     make(map[types.Id]observable.Observable[float32]),
		config:   config,
		limbMode: limbMode,
	}, nil
}

// TODO implement fuel
func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()

	p.fuel[carId] = observable.Create(float32(p.carConfig[carId].FuelConfig.TankSize))
	p.fuelState[carId] = fuelState{true, 0}

	car.Controller().TriggerValue().RegisterObserver(func(v uint8, annotations observable.Annotations) {
		if v > 0 {
			liter := p.fuel[carId].Get()
			if liter > 0 {
				s := p.fuelState[carId]
				used := calculateFuelState(p.carConfig[carId].FuelConfig.BurnRate, p.fuelState[carId].consumed, v)
				if float32(p.carConfig[carId].FuelConfig.TankSize) >= used {
					s.consumed = used
				}
				p.fuel[carId].Publish()
			}
		}
		panic("Implement me")
	})

	car.Deslotted().RegisterObserver(func(b bool, annotations observable.Annotations) {
		panic("Implement me")

	})

	car.Pit().RegisterObserver(func(b bool, a observable.Annotations) {
		panic("Implement me")

		// TODO Clear consumption when entering pit
	})

	p.fuel[carId].RegisterObserver(func(f float32, a observable.Annotations) {
		panic("Implement me")

	})

	p.fuel[carId].Modifier(func(f float32) (float32, bool) {
		m := f
		if p.fuelState[carId].enabled {
			m = f - p.fuelState[carId].consumed
		}
		return m, p.fuelState[carId].enabled
	}, 1)

}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.RaceStatus, a observable.Annotations) {
		p.status = status
	})
}

func calculateFuelState(burnRate types.LiterPerSecond, fuel float32, triggerValue uint8) float32 {
	used := (float32(triggerValue) / 100) * float32(burnRate)
	return used + fuel
}
