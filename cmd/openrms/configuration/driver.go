package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/implement/generator"
	"github.com/qvistgaard/openrms/internal/implement/oxigen"
)

type DriverConfiguration struct {
	Implement struct {
		Plugin string
	}
}

// Driver initializes and returns an implementation of the `implement.Implementer` interface based on the provided configuration.
// It takes a `Config` map, which should contain the driver configuration, and an optional pointer to a driver name.
//
// The `conf` parameter should be a `Config` map containing the driver configuration settings.
//
// The `driver` parameter is an optional pointer to a string that represents the name of the driver to be used.
// If `driver` is provided and not empty, it will take precedence over the driver name specified in the configuration.
//
// Example usage:
//
//   // Load driver configuration from a previously loaded configuration map.
//   driverConfig := configMap["driver"].(map[string]interface{})
//
//   // Specify the driver name (optional).
//   driverName := "oxigen"
//
//   // Initialize the driver based on the configuration.
//   driver, err := Driver(driverConfig, &driverName)
//   if err != nil {
//       log.Fatal("Failed to initialize the driver: ", err)
//   }
//
//   // Use the 'driver' as an implementation of the 'implement.Implementer' interface.
//
// Returns:
//   - An implementation of the 'implement.Implementer' interface corresponding to the specified driver.
//   - An error if there was an issue initializing the driver or if the specified driver is unknown.
func Driver(conf Config, driver *string) (implement.Implementer, error) {
	c := &DriverConfiguration{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read driver configuration")
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
			return nil, errors.WithMessage(err, "failed to read oxigen driver configuration")
		}
		return oxigen.New(*c)
	case "generator":
		c := &generator.Config{}
		err := mapstructure.Decode(conf, c)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read generator driver configuration")
		}

		return generator.New(*c)
	default:
		return nil, errors.New("Unknown implementer: " + c.Implement.Plugin)
	}
}
