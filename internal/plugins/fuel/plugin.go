// Package fuel provides a plugin for managing fuel levels in cars during a race.
package fuel

import (
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
)

// Plugin represents the fuel management plugin.
// It provides functionality for monitoring and managing fuel levels in cars during a race.
type Plugin struct {
	config    Config
	carConfig map[types.Id]CarSettings
	state     map[types.Id]*state
	status    race.RaceStatus
	limbMode  *limbmode.Plugin
}

// state represents the fuel state of an individual car.
type state struct {
	enabled  bool
	consumed float32
	machine  *stateless.StateMachine
	fuel     observable.Observable[float32]
}

// New creates a new instance of the fuel plugin.
func New(config Config, limbMode *limbmode.Plugin) (*Plugin, error) {
	return &Plugin{
		config:   config,
		limbMode: limbMode,
		state:    make(map[types.Id]*state),
	}, nil
}

// ConfigureCar configures the fuel management for a car.
// It sets up the state machine and observers to monitor and update the fuel level of the car.
// The function performs the following steps:
// 1. Retrieves the car's fuel configuration or uses default values if not configured.
// 2. Creates a state machine for managing fuel updates.
//   - This state machine controls when and how the fuel level is updated based on trigger events.
//   - It also transitions the car's fuel state between on track and deslotted.
//   - The machine handles events triggered by changes in the car's trigger value and deslotting status.
//
// 3. Registers observers for the car's trigger value, deslotting status, and pit status.
//   - The trigger value observer triggers fuel updates based on changes in the trigger value.
//   - The deslotting observer transitions the fuel state when the car is deslotted.
//   - The pit observer resets the consumed fuel and sets the fuel level to the tank's capacity when the car enters the pit.
//
// 4. Registers a modifier function to update the fuel level based on consumption with a priority of 1.
// 5. Registers an observer to check if the fuel level drops to or below zero, signaling a limb mode activation.
//
// Parameters:
//   - car: The car to configure fuel management for.
func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()
	p.state[carId] = &state{}
	carState := p.state[carId]
	config := p.carConfig[carId].FuelConfig
	if config == nil {
		config = p.config.Car.Defaults.FuelConfig
	}

	carState.machine = machine(handleUpdateFuelLevel(carState, config.TankSize, config.BurnRate))
	carState.fuel = observable.Create(float32(config.TankSize))

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

	// Register a modifier function to update the fuel level based on consumption.
	// The modifier function subtracts the consumed fuel from the current fuel level.
	// It has a priority of 1, ensuring it is applied before other modifiers.
	carState.fuel.Modifier(func(f float32) (float32, bool) {
		return f - carState.consumed, true
	}, 1)

	// Register an observer to check if the fuel level drops to or below zero.
	// When this happens, it signals a limb mode activation for the car.
	carState.fuel.RegisterObserver(func(f float32, a observable.Annotations) {
		if f <= 0 {
			p.limbMode.LimbMode(carId).Set(true)
		}
	})

}

func (p *Plugin) InitializeCar(c *car.Car) {

}

// Priority returns the priority of the fuel plugin.
func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "fuel"
}

// Fuel returns the observable fuel level for a given car.
func (p *Plugin) Fuel(car types.Id) observable.Observable[float32] {
	return p.state[car].fuel
}

// ConfigureRace configures the fuel plugin for a race.
// It registers an observer for monitoring the race status.
func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.RaceStatus, a observable.Annotations) {
		p.status = status
	})
}
