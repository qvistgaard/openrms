package configuration

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/state/track"
)

// Track initializes and returns a new track instance based on the provided configuration and drivers implementation.
// It takes a `Config` map containing the track configuration settings and an implementation of the `drivers.Driver` interface.
//
// The `conf` parameter should be a `Config` map containing the track configuration settings.
//
// The `drivers` parameter is an implementation of the `drivers.Driver` interface that will be used to manage the track.
//
// Example usage:
//
//	// Load track configuration from a previously loaded configuration map.
//	trackConfig := configMap["track"].(map[string]interface{})
//
//	// Initialize the drivers (drivers initialization not shown here).
//	var drivers drivers.Driver
//
//	// Initialize the track instance based on the configuration and drivers.
//	trackInstance, err := Track(trackConfig, drivers)
//	if err != nil {
//	    log.Fatal("Failed to initialize the track: ", err)
//	}
//
//	// Use the 'trackInstance' for managing track-related operations.
//
// Returns:
//   - A new instance of the 'track.Track' type representing the initialized track.
//   - An error if there was an issue initializing the track instance.
func Track(conf Config, driver drivers.Driver) (*track.Track, error) {
	c := &track.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read track configuration")
	}

	return track.New(*c, driver)
}
