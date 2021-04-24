package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

type Config struct {
	Car struct {
		Cars []map[string]interface{}
	}
}

type CarConfigRepository struct {
	cars map[state.CarId]map[string]interface{}
}

func (c *CarConfigRepository) GetCarById(id state.CarId) map[string]interface{} {
	if settings, ok := c.cars[id]; ok {
		return settings
	}
	return make(map[string]interface{})
}

func CreateFromConfig(ctx *context.Context) (*CarConfigRepository, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	ccr := new(CarConfigRepository)
	ccr.cars = make(map[state.CarId]map[string]interface{})
	for _, cs := range c.Car.Cars {
		if id, ok := cs["id"]; ok {
			i := id.(int)
			ccr.cars[state.CarId(i)] = cs
		}
	}
	return ccr, nil
}
