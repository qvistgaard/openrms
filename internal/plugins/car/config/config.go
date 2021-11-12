package config

import (
	ctx "context"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state/rx/car"
	"github.com/qvistgaard/openrms/internal/types"
)

type Config struct {
	Car struct {
		Defaults map[string]interface{}
		Cars     []CarConfig
	}
}
type CarConfig struct {
	Id       types.Id      `yaml:"id"`
	MaxSpeed types.Percent `mapstructure:"max-speed"`
}

type CarConfigRepository struct {
	cars     map[types.Id]*car.Car
	config   map[types.Id]CarConfig
	context  *application.Context
	defaults map[string]interface{}
}

func CreateFromConfig(ctx *application.Context) (*CarConfigRepository, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	ccr := new(CarConfigRepository)
	ccr.cars = make(map[types.Id]*car.Car)
	ccr.config = make(map[types.Id]CarConfig)
	ccr.defaults = c.Car.Defaults
	ccr.context = ctx

	for _, cs := range c.Car.Cars {
		ccr.config[cs.Id] = cs
	}
	return ccr, nil
}

func (c *CarConfigRepository) Get(id types.Id, ctx ctx.Context) (*car.Car, bool, bool) {
	carCreated := false
	if _, ok := c.cars[id]; !ok {
		if _, ok := c.config[id]; !ok {
			c.config[id] = CarConfig{}
		}

		// i := merge.Merge(c.defaults, c.config[id])
		c.cars[id] = car.NewCar(c.context.Implement, id)

		for _, r := range c.context.Rules.CarRules() {
			r.InitializeCarState(c.cars[id])
		}
		c.cars[id].Init(ctx, c.context.Postprocessors.ValuePostProcessor())

		// c.cars[id].MaxSpeed().Set(c.config[id].MaxSpeed)
		// c.cars[id].Rx.Percent().Set(c.cars[id].GetSettings("max-speed").(state.Speed))
		carCreated = true
	}
	return c.cars[id], true, carCreated
}

func (c *CarConfigRepository) Exists(id types.Id) bool {
	_, ok := c.cars[id]
	return ok
}

func (c *CarConfigRepository) All() []*car.Car {
	cars := make([]*car.Car, 0, len(c.cars))
	for _, car := range c.cars {
		cars = append(cars, car)
	}
	return cars
}
