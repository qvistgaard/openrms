// Package fuel provides a plugin for managing fuel levels in cars during a race.
package fuel

import (
	"embed"
	"errors"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/plugins/commentary"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/utils"
	log "github.com/sirupsen/logrus"
)

//go:embed commentary/out_of_fuel.txt
var announcements embed.FS

// Plugin represents the fuel management plugin.
// It provides functionality for monitoring and managing fuel levels in cars during a race.
type Plugin struct {
	config     Config
	carConfig  map[types.CarId]CarSettings
	state      map[types.CarId]*state
	usage      map[types.CarId][]float32
	status     race.Status
	limbMode   *limbmode.Plugin
	commentary *commentary.Plugin
}

// state represents the fuel state of an individual car.
type state struct {
	enabled bool
	// consumed is current amount of fuel consumed.
	consumed  float32
	machine   *stateless.StateMachine
	fuel      observable.Observable[float32]
	config    FuelConfig
	average   average
	announced bool
}

// New creates a new instance of the fuel plugin.
func New(config Config, limbMode *limbmode.Plugin, commentary *commentary.Plugin) (*Plugin, error) {
	return &Plugin{
		config:     config,
		limbMode:   limbMode,
		commentary: commentary,
		state:      make(map[types.CarId]*state),
		usage:      make(map[types.CarId][]float32),
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
	carState.config = *config

	carState.machine = machine(handleUpdateFuelLevel(carState, config.TankSize, config.BurnRate))
	carState.fuel = createFuelObserver(float32(carState.config.TankSize), fuelModifier(carState), limpModeObserver(carId, p.limbMode))

	car.LastLap().RegisterObserver(func(lap types.Lap) {
		p.usage[carId] = append(p.usage[carId], p.state[carId].consumed)
		carState.average = carState.average.reportUsage(carState.consumed)

		f := carState.fuel.Get() / carState.average.average
		if f < 5 && p.config.Plugin.Fuel.Commentary {
			line, err := utils.RandomLine(announcements, "commentary/out_of_fuel.txt")
			if err == nil && !carState.announced {
				template, _ := utils.ProcessTemplate(line, car.TemplateData())
				p.commentary.Announce(template)
				carState.announced = true
			}
		}

	})

	car.Controller().TriggerValue().RegisterObserver(func(v uint8) {
		if p.status == race.Running {
			err := carState.machine.Fire(triggerUpdateFuelLevel, v)
			if err != nil {
				log.Error(err)
			}
		}
	})

	car.Deslotted().RegisterObserver(func(b bool) {
		if b {
			err := carState.machine.Fire(triggerCarDeslotted)
			if err != nil {
				log.Error(err)
			}
		} else {
			err := carState.machine.Fire(triggerCarOnTrack)
			if err != nil {
				log.Error(err)
			}
		}
	})
}

func (p *Plugin) InitializeCar(_ *car.Car) {

}

// Priority returns the priority of the fuel plugin.
func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "fuel"
}

// Fuel returns the observable fuel level for a given car.
func (p *Plugin) Fuel(car types.CarId) (observable.Observable[float32], error) {
	if f, ok := p.state[car]; ok {
		return f.fuel, nil
	}
	return nil, errors.New("car not found")
}

// createFuelObserver initializes and returns an observable fuel level value for a car,
// with mechanisms to modify and observe the fuel level based on defined rules.
//
// Parameters:
//   - tankSize: The initial size of the fuel tank, representing the starting fuel level.
//   - modifier: A function that adjusts the fuel level based on certain criteria (e.g., consumption).
//     This function should return the modified fuel level and a boolean indicating success.
//   - limbmodeObserver: A function that observes changes in fuel level and performs actions based on
//     those changes, such as activating limb mode if the fuel level drops too low.
//
// Returns:
//   - An observable.Value[float32] instance representing the car's fuel level. This object allows
//     registration of modifier and observer functions to dynamically manage fuel levels.
//
// The function uses an observable value to manage the fuel level, enabling real-time adjustments
// and monitoring. Modifier functions are applied to adjust the fuel level as the simulation or
// operation progresses, while observer functions can trigger additional actions based on fuel level changes.
func createFuelObserver(tankSize float32, modifier func(f float32) (float32, bool), limbmodeObserver func(f float32)) *observable.Value[float32] {
	fuel := observable.Create(tankSize)

	// Register a modifier function to update the fuel level based on consumption.
	// The modifier function subtracts the consumed fuel from the current fuel level.
	// It has a priority of 1, ensuring it is applied before other modifiers.
	fuel.Modifier(modifier, 1)

	// Register an observer to check if the fuel level drops to or below zero.
	// When this happens, it signals a limb mode activation for the car.
	fuel.RegisterObserver(limbmodeObserver)

	return fuel
}

// fuelModifier constructs a fuel level adjustment function based on a car's fuel consumption state.
// This function is intended to be used as a modifier with the createFuelObserver function, allowing
// dynamic adjustment of the fuel level during the simulation or operation of the car.
//
// Parameters:
//   - carState: A pointer to the `state` struct, which encapsulates the car's fuel consumption data,
//     including the amount of fuel already consumed.
//
// Returns:
//   - A closure that takes the current fuel level (f float32) and returns the adjusted fuel level
//     (after accounting for consumption) and a boolean indicating successful adjustment. This closure
//     matches the signature expected by createFuelObserver's modifier parameter.
//
// The returned modifier function deducts the consumed fuel from the given fuel level, aiding in
// the simulation of fuel consumption over time.
func fuelModifier(carState *state) func(f float32) (float32, bool) {
	return func(f float32) (float32, bool) {
		return f - carState.consumed, true
	}
}

// limpModeObserver creates a function to monitor the car's fuel level and activate limp mode if the fuel
// level falls to or below zero. This observer function is designed for use with createFuelObserver,
// allowing it to be registered as an observer for fuel level changes.
//
// Parameters:
// - carId: The unique identifier for the car, used to determine which car's limp mode should be activated.
// - p: A pointer to a limbmode.Plugin instance, enabling the function to set the car into limp mode.
//
// Returns:
//   - A function that takes the current fuel level (f float32) and activates limp mode for the car
//     identified by carId if the fuel level is zero or less. This function matches the signature expected
//     by createFuelObserver's RegisterObserver method.
//
// When the fuel level drops to 0 or below, indicating the car has run out of fuel, this observer
// triggers the limp mode to preserve the car's operational integrity.
func limpModeObserver(carId types.CarId, p *limbmode.Plugin) func(f float32) {
	return func(f float32) {
		if f <= 0 {
			p.LimbMode(carId).Set(true)
		}
	}
}

// ConfigureRace sets up race-related configurations for the plugin, particularly focusing on
// race status observation. It registers an observer function to monitor changes in the race's status.
//
// Parameters:
// - r: A pointer to a `race.Race` instance, representing the race to be configured with this plugin.
//
// This method leverages the race's observable status feature to keep the plugin informed about
// the race's current state. Upon any change in race status, the registered observer function
// (raceStatusObserver) is invoked, allowing the plugin to respond appropriately to events such as
// the race stopping.
//
// Usage:
// Intended to be called during the race setup phase, this method equips the plugin with the
// ability to react dynamically to changes in the race's lifecycle, ensuring relevant states
// within the plugin are reset or updated in response to race events.
func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(raceStatusObserver(p))
}

// raceStatusObserver creates and returns a closure that acts as an observer for race status changes.
// This observer is responsible for updating the plugin's internal state based on the current status
// of the race, particularly resetting consumed fuel states upon race completion.
//
// Parameters:
// - p: A pointer to a `Plugin` instance, which holds the state to be updated by this observer.
//
// Returns:
//   - A function that takes a race.Status and performs actions based on that status. Specifically,
//     when the race status changes to `race.Stopped`, this function resets the consumed fuel values
//     for all cars managed by the plugin and triggers an update to reflect these resets.
//
// The observer function directly manipulates the internal state of the plugin in response to
// race lifecycle events, ensuring that the plugin accurately reflects the current state of
// the race and is ready for a new race start with reset fuel consumption values.
//
// Usage:
// This function is not called directly but is passed as an argument to the race status
// observable's RegisterObserver method, allowing it to be automatically invoked whenever
// the race status changes.
func raceStatusObserver(p *Plugin) func(status race.Status) {
	return func(status race.Status) {
		p.status = status
		if status == race.Stopped {
			for _, s := range p.state {
				s.consumed = 0
				s.fuel.Update()
			}
		}
	}
}

func (p *Plugin) ConfigurePitSequence(carId types.CarId) pit.Sequence {
	return NewSequence(p.state[carId])
}
