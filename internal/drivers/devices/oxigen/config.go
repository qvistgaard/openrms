package oxigen

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/drivers"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string
		}
	}
}

func New(config Config) (drivers.Driver, error) {
	connection, err := CreateUSBConnection(&config.Implement.Oxigen.Port)
	if err != nil {
		return nil, errors.New("Failed to open serial to USB Device (" + config.Implement.Oxigen.Port + "): " + err.Error())
	}
	return CreateImplement(connection)
}
