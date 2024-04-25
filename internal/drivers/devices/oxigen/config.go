package oxigen

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/devices/oxigen/serial"
	"github.com/rs/zerolog"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string
		}
	}
}

func New(logger zerolog.Logger, config Config) (drivers.Driver, error) {
	connection, err := serial.CreateUSBConnection(logger, &config.Implement.Oxigen.Port)
	if err != nil {
		return nil, errors.New("Failed to open serial to USB Device (" + config.Implement.Oxigen.Port + "): " + err.Error())
	}
	return CreateImplement(logger, connection)
}
