package fuel

import (
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
)

type Plugin struct {
	config    Config
	carConfig map[types.Id]CarSettings
	state     map[types.Id]*state
	status    race.RaceStatus
	limbMode  *limbmode.Plugin
}

type state struct {
	enabled  bool
	consumed float32
	machine  *stateless.StateMachine
	fuel     observable.Observable[float32]
}

func New(config Config, limbMode *limbmode.Plugin) (*Plugin, error) {
	return &Plugin{
		// fuel:     make(map[types.Id]observable.Observable[float32]),
		config:   config,
		limbMode: limbMode,
		state:    make(map[types.Id]*state),
	}, nil
}

// TODO implement fuel
func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()
	p.state[carId] = &state{}
	carState := p.state[carId]
	config := p.carConfig[carId].FuelConfig
	if config == nil {
		config = &FuelConfig{
			TankSize:     80,
			StartingFuel: 60,
			BurnRate:     100,
			FlowRate:     0,
		}
	}

	carState.machine = machine(handleUpdateFuelLevel(carState, config.TankSize, config.BurnRate))
	carState.fuel = observable.Create(float32(config.TankSize))
	// .fuelState[carId] = fuelState{true, 0}

	car.Controller().TriggerValue().RegisterObserver(func(v uint8, annotations observable.Annotations) {
		carState.machine.Fire(triggerUpdateFuelLevel, v)
	})

	car.Deslotted().RegisterObserver(func(b bool, annotations observable.Annotations) {
		if b {
			carState.machine.Fire(triggerCarDeslotted)
		} else {
			carState.machine.Fire(triggerCarOnTrack)
		}
	})

	car.Pit().RegisterObserver(func(b bool, a observable.Annotations) {
		carState.consumed = 0
		carState.fuel.Set(float32(config.TankSize))
	})

	carState.fuel.Modifier(func(f float32) (float32, bool) {
		return f - carState.consumed, true
	}, 1)

	carState.fuel.RegisterObserver(func(f float32, a observable.Annotations) {
		if f <= 0 {
			p.limbMode.LimbMode(carId).Set(true)
		}
	})

}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Fuel(car types.Id) observable.Observable[float32] {
	return p.state[car].fuel
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.RaceStatus, a observable.Annotations) {
		p.status = status
	})
}

func calculateFuelState(burnRate float32, consumed float32, triggerValue uint8) float32 {
	return ((float32(triggerValue) / 100) * burnRate) + consumed
}
