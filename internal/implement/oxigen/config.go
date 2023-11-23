package oxigen

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/implement"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string
		}
	}
}

func New(config Config) (implement.Implementer, error) {
	connection, err := CreateUSBConnection(config.Implement.Oxigen.Port)
	if err != nil {
		return nil, errors.New("Failed to open connection to USB Device (" + config.Implement.Oxigen.Port + "): " + err.Error())
	}
	return CreateImplement(connection)
}
