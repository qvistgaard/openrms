package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/state/race"
)

// Race initializes and returns a new race instance based on the provided configuration and drivers implementation.
// It takes a `Config` map containing the race configuration settings and an implementation of the `drivers.Driver` interface.
//
// The `conf` parameter should be a `Config` map containing the race configuration settings.
//
// The `drivers` parameter is an implementation of the `drivers.Driver` interface that will be used to manage the race.
//
// Example usage:
//
//	// Load race configuration from a previously loaded configuration map.
//	raceConfig := configMap["race"].(map[string]interface{})
//
//	// Initialize the drivers (drivers initialization not shown here).
//	var drivers drivers.Driver
//
//	// Initialize the race instance based on the configuration and drivers.
//	raceInstance, err := Race(raceConfig, drivers)
//	if err != nil {
//	    log.Fatal("Failed to initialize the race: ", err)
//	}
//
//	// Use the 'raceInstance' for managing race-related operations.
//
// Returns:
//   - A new instance of the 'race.Race' type representing the initialized race.
//   - An error if there was an issue initializing the race instance.
func Race(conf Config, driver drivers.Driver) (*race.Race, error) {
	c := &race.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read race configuration")
	}

	return race.New(*c, driver)
}
