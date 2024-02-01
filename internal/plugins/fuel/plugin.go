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
	enabled   bool
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
	carState.fuel = observable.Create(float32(config.TankSize))

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

	car.Pit().RegisterObserver(func(b bool) {
		carState.consumed = 0
		carState.announced = false
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
	carState.fuel.RegisterObserver(func(f float32) {
		if f <= 0 {
			p.limbMode.LimbMode(carId).Set(true)
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

// ConfigureRace configures the fuel plugin for a race.
// It registers an observer for monitoring the race status.
func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.Status) {
		p.status = status
		if status == race.Stopped {
			for _, s := range p.state {
				s.consumed = 0
				s.fuel.Update()
			}
		}
	})
}

func (p *Plugin) ConfigurePitSequence(carId types.CarId) pit.Sequence {
	return NewSequence(p.state[carId])
}
