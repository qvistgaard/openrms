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
	cars    map[state.CarId]*state.Car
	config  map[state.CarId]map[string]interface{}
	context *context.Context
}

func CreateFromConfig(ctx *context.Context) (*CarConfigRepository, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	ccr := new(CarConfigRepository)
	ccr.cars = make(map[state.CarId]*state.Car)
	ccr.config = make(map[state.CarId]map[string]interface{})
	ccr.context = ctx
	for _, cs := range c.Car.Cars {
		if id, ok := cs["id"]; ok {
			i := id.(int)
			ccr.config[state.CarId(i)] = cs
		}
	}
	return ccr, nil
}

func (c *CarConfigRepository) Get(id state.CarId) (*state.Car, bool) {
	if _, ok := c.cars[id]; !ok {
		if _, ok := c.config[id]; !ok {
			c.config[id] = make(map[string]interface{})
		}
		c.cars[id] = state.CreateCar(id, c.config[id], c.context.Rules)
	}
	return c.cars[id], true
}

func (c *CarConfigRepository) All() []*state.Car {
	panic("implement me")
}
