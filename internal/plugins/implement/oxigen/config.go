package oxigen

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string
		}
	}
}

func CreateFromConfig(context *context.Context) (*Oxigen, error) {
	c := &Config{}
	err := mapstructure.Decode(context.Config, c)
	if err != nil {
		return nil, err
	}

	connection, err := CreateUSBConnection(c.Implement.Oxigen.Port)
	if err != nil {
		return nil, errors.New("Failed to open connection to USB Device (" + c.Implement.Oxigen.Port + "): " + err.Error())
	}
	return CreateImplement(connection)
}
