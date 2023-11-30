package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/plugins"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/car/repository"
)

// CarRepository creates a new car repository with the provided configuration and plugins.
// It takes a `Config` struct and a list of plugins as input and returns a `Repository`
// interface representing a car repository and an error, if any.
//
// The `Config` struct contains the configuration settings required to initialize the car repository.
//
// The `plugins` parameter is a list of plugins that can be used to extend the functionality
// of the car repository.
//
// Example usage:
//
//	conf := configuration.Config{
//	    // Populate the configuration settings as needed.
//	}
//
//	plugins := plugins.List{
//	    // Populate the list of plugins as needed.
//	}
//
//	carRepo, err := CarRepository(conf, plugins)
//	if err != nil {
//	    log.Fatal("Failed to create car repository: ", err)
//	}
//
//	// Use the car repository for managing car-related data.
//
// Returns:
//   - A `Repository` interface representing the car repository.
//   - An error if there was an issue creating the repository.
func CarRepository(conf Config, driver drivers.Driver, plugins plugins.List) (repository.Repository, error) {
	c := &car.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read car configuration")
	}

	return repository.New(*c, driver, plugins), err
}
