package fuel

import (
	"github.com/qvistgaard/openrms/internal/types"
)

// Config represents the configuration for managing car settings including fuel configuration.
type Config struct {
	Car struct {
		// Defaults represents the default car settings that can be inherited by other cars.
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
}

// CarSettings represents the configuration settings for an individual car.
type CarSettings struct {
	Id         *types.Id   `mapstructure:"id"`
	FuelConfig *FuelConfig `mapstructure:"fuel"`
}

// FuelConfig represents the fuel-related configuration settings for a car.
type FuelConfig struct {
	// TankSize is the size of the fuel tank in liters.
	TankSize uint8 `mapstructure:"tank-size"`

	// StartingFuel is the initial fuel level when the car is placed on the track.
	StartingFuel uint8 `mapstructure:"starting-fuel"`

	// BurnRate is the rate at which the car consumes fuel per unit of time.
	BurnRate float32 `mapstructure:"burn-rate"`

	// FlowRate is the rate at which fuel flows into the car's fuel tank when refuelling.
	FlowRate float32 `mapstructure:"flow-rate"`
}
