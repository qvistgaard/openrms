package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/plugins/implement/generator"
	"github.com/qvistgaard/openrms/internal/plugins/implement/oxigen"
)

type ImplementConfig struct {
	Implement struct {
		Plugin string
	}
}

func CreateImplement(context *application.Context) error {
	c := &ImplementConfig{}
	err := mapstructure.Decode(context.Config, c)
	if err != nil {
		return err
	}

	switch c.Implement.Plugin {
	case "oxigen":
		context.Implement, err = oxigen.CreateFromConfig(context)
	case "generator":
		// TODO: recreate generating implement
		context.Implement, err = generator.CreateFromConfig(context)
	default:
		return errors.New("Unknown implementer: " + c.Implement.Plugin)
	}

	return err
}
