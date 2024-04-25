package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/devices/generator"
	"github.com/qvistgaard/openrms/internal/drivers/devices/oxigen"
	"github.com/rs/zerolog"
)

type DriverConfiguration struct {
	Implement struct {
		Plugin string
	}
}

// Driver initializes and returns an implementation of the `drivers.Driver` interface based on the provided configuration.
// It takes a `Config` map, which should contain the drivers configuration, and an optional pointer to a drivers name.
//
// The `conf` parameter should be a `Config` map containing the drivers configuration settings.
//
// The `drivers` parameter is an optional pointer to a string that represents the name of the drivers to be used.
// If `drivers` is provided and not empty, it will take precedence over the drivers name specified in the configuration.
//
// Example usage:
//
//	// Load drivers configuration from a previously loaded configuration map.
//	driverConfig := configMap["drivers"].(map[string]interface{})
//
//	// Specify the drivers name (optional).
//	driverName := "oxigen"
//
//	// Initialize the drivers based on the configuration.
//	drivers, err := Driver3x(driverConfig, &driverName)
//	if err != nil {
//	    log.Fatal("Failed to initialize the drivers: ", err)
//	}
//
//	// Use the 'drivers' as an implementation of the 'drivers.Driver' interface.
//
// Returns:
//   - An implementation of the 'drivers.Driver' interface corresponding to the specified drivers.
//   - An error if there was an issue initializing the drivers or if the specified drivers is unknown.
func Driver(logger zerolog.Logger, conf Config, driver *string) (drivers.Driver, error) {
	c := &DriverConfiguration{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read drivers configuration")
	}

	var plugin string
	if driver == nil || *driver == "" {
		plugin = c.Implement.Plugin
	} else {
		plugin = *driver
	}

	switch plugin {
	case "oxigen":
		c := &oxigen.Config{}
		err := mapstructure.Decode(conf, c)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read oxigen drivers configuration")
		}
		return oxigen.New(logger, *c)
	case "generator":
		c := &generator.Config{}
		err := mapstructure.Decode(conf, c)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read generator drivers configuration")
		}

		return generator.New(*c)
	default:
		return nil, errors.New("Unknown implementer: " + c.Implement.Plugin)
	}
}
