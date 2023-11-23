package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/race"
)

// FuelPlugin initializes and returns a new fuel plugin instance based on the provided configuration and LimpMode plugin.
// It takes a `Config` map containing the fuel plugin configuration settings and an optional `limbmode.Plugin` instance.
//
// The `conf` parameter should be a `Config` map containing the fuel plugin configuration settings.
//
// The `limpMode` parameter is an optional instance of the `limbmode.Plugin` type that can be provided if needed.
// If not needed, you can pass `nil` for this parameter.
//
// Example usage:
//
//   // Load fuel plugin configuration from a previously loaded configuration map.
//   fuelConfig := configMap["fuel"].(map[string]interface{})
//
//   // Initialize the LimpMode plugin (LimpMode plugin initialization not shown here).
//   var limpMode *limbmode.Plugin
//
//   // Initialize the fuel plugin instance based on the configuration and LimpMode plugin.
//   fuelInstance, err := FuelPlugin(fuelConfig, limpMode)
//   if err != nil {
//       log.Fatal("Failed to initialize the fuel plugin: ", err)
//   }
//
//   // Use the 'fuelInstance' for managing fuel-related operations.
//
// Returns:
//   - A new instance of the 'fuel.Plugin' type representing the initialized fuel plugin.
//   - An error if there was an issue initializing the fuel plugin instance.
func FuelPlugin(conf Config, limpMode *limbmode.Plugin) (*fuel.Plugin, error) {
	c := &fuel.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}

	return fuel.New(*c, limpMode)
}

// RacePlugin initializes and returns a new race plugin instance based on the provided configuration.
// It takes a `Config` map containing the race plugin configuration settings.
//
// The `conf` parameter should be a `Config` map containing the race plugin configuration settings.
//
// Example usage:
//
//   // Load race plugin configuration from a previously loaded configuration map.
//   raceConfig := configMap["race"].(map[string]interface{})
//
//   // Initialize the race plugin instance based on the configuration.
//   raceInstance, err := RacePlugin(raceConfig)
//   if err != nil {
//       log.Fatal("Failed to initialize the race plugin: ", err)
//   }
//
//   // Use the 'raceInstance' for managing race-related operations.
//
// Returns:
//   - A new instance of the 'race.Plugin' type representing the initialized race plugin.
//   - An error if there was an issue initializing the race plugin instance.
func RacePlugin(_ Config) (*race.Plugin, error) {
	return race.New()
}

// LimbModePlugin initializes and returns a new LimpMode plugin instance based on the provided configuration.
// It takes a `Config` map containing the LimpMode plugin configuration settings.
//
// The `conf` parameter should be a `Config` map containing the LimpMode plugin configuration settings.
//
// Example usage:
//
//   // Load LimpMode plugin configuration from a previously loaded configuration map.
//   limpModeConfig := configMap["limpmode"].(map[string]interface{})
//
//   // Initialize the LimpMode plugin instance based on the configuration.
//   limpModeInstance, err := LimbModePlugin(limpModeConfig)
//   if err != nil {
//       log.Fatal("Failed to initialize the LimpMode plugin: ", err)
//   }
//
//   // Use the 'limpModeInstance' for managing LimpMode-related operations.
//
// Returns:
//   - A new instance of the 'limbmode.Plugin' type representing the initialized LimpMode plugin.
//   - An error if there was an issue initializing the LimpMode plugin instance.
func LimbModePlugin(_ Config) (*limbmode.Plugin, error) {
	return &limbmode.Plugin{}, nil
}
