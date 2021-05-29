package config

import (
	"github.com/divideandconquer/go-merge/merge"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

type Config struct {
	Car struct {
		Defaults map[string]interface{}
		Cars     []map[string]interface{}
	}
}

type CarConfigRepository struct {
	cars     map[state.CarId]*state.Car
	config   map[state.CarId]map[string]interface{}
	context  *context.Context
	defaults map[string]interface{}
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
	ccr.defaults = c.Car.Defaults
	ccr.context = ctx
	for _, cs := range c.Car.Cars {
		if id, ok := cs["id"]; ok {
			i := id.(int)
			ccr.config[state.CarId(i)] = cs
		}
	}
	return ccr, nil
}

func (c *CarConfigRepository) Get(id state.CarId) (*state.Car, bool, bool) {
	carCreated := false
	if _, ok := c.cars[id]; !ok {
		if _, ok := c.config[id]; !ok {
			c.config[id] = make(map[string]interface{})
		}

		i := merge.Merge(c.defaults, c.config[id])
		c.cars[id] = state.CreateCar(id, i.(map[string]interface{}), c.context.Rules)
		carCreated = true
	}
	return c.cars[id], true, carCreated
}

func (c *CarConfigRepository) Exists(id state.CarId) bool {
	_, ok := c.cars[id]
	return ok
}

func (c *CarConfigRepository) All() []*state.Car {
	cars := make([]*state.Car, 0, len(c.cars))
	for _, car := range c.cars {
		cars = append(cars, car)
	}
	return cars
}
