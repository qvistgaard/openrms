package oxigen

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string `yaml:"port"`
		} `yaml:"oxigen"`
	} `yaml:"implement"`
}

func CreateFromConfig(context *context.Context) (*Oxigen, error) {
	c := &Config{}
	err := mapstructure.Decode(context, c)
	if err != nil {
		return nil, err
	}

	connection, err := CreateUSBConnection(c.Implement.Oxigen.Port)
	if err != nil {
		return nil, errors.New("Failed to open connection to USB Device: " + err.Error())
	}
	return CreateImplement(connection)
}
