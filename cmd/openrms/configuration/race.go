package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/race"
)

// Race initializes and returns a new race instance based on the provided configuration and driver implementation.
// It takes a `Config` map containing the race configuration settings and an implementation of the `implement.Implementer` interface.
//
// The `conf` parameter should be a `Config` map containing the race configuration settings.
//
// The `driver` parameter is an implementation of the `implement.Implementer` interface that will be used to manage the race.
//
// Example usage:
//
//   // Load race configuration from a previously loaded configuration map.
//   raceConfig := configMap["race"].(map[string]interface{})
//
//   // Initialize the driver (driver initialization not shown here).
//   var driver implement.Implementer
//
//   // Initialize the race instance based on the configuration and driver.
//   raceInstance, err := Race(raceConfig, driver)
//   if err != nil {
//       log.Fatal("Failed to initialize the race: ", err)
//   }
//
//   // Use the 'raceInstance' for managing race-related operations.
//
// Returns:
//   - A new instance of the 'race.Race' type representing the initialized race.
//   - An error if there was an issue initializing the race instance.
func Race(conf Config, driver implement.Implementer) (*race.Race, error) {
	c := &race.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read race configuration")
	}

	return race.New(*c, driver)
}
