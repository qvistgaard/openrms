package plugins

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/race"
)

type Plugin interface {
	Priority() int
	Name() string
}

type Car interface {

	// ConfigureCar is responsible for setting up the plugin-specific
	// configuration for a Car object. This method is invoked when a new car
	// is detected, allowing each plugin to apply its unique configuration to the
	// car's state. This method should not alter the car itself but should
	// prepare any plugin-specific settings or prerequisites that are necessary
	// before the plugin is initialized.
	//
	// Parameters:
	//   car - A pointer to the Car object representing the state of the car.
	//
	// Note: This method should focus on preparing the plugin's configuration
	// and should not depend on the initialization of the plugin, as it precedes
	// the InitializeCarState method.
	ConfigureCar(*car.Car)

	// InitializeCar init
	// Deprecated: seems to be unused
	InitializeCar(*car.Car)
}

type Race interface {
	// ConfigureRace
	//
	// Deprecated: no longer needed
	ConfigureRace(*race.Race)
}

type List interface {
	Car() []Car
	Race() []Race
	Append(Plugin) Plugin
}
