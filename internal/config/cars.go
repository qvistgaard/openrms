package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/plugins/car/config"
)

type CarsConfig struct {
	Car struct {
		Plugin string
	}
}

func CreateCarRepository(ctx *context.Context) error {
	c := &CarsConfig{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return err
	}

	switch c.Car.Plugin {
	case "config":
		ctx.Cars, err = config.CreateFromConfig(ctx)
	default:
		return errors.New("no car configuration found")
	}
	return err
}
