package config

import (
	ctx "context"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state/car"
	config "github.com/qvistgaard/openrms/internal/state/config/car"
	"github.com/qvistgaard/openrms/internal/types"
)

type CarConfigRepository struct {
	cars     map[types.Id]*car.Car
	config   map[types.Id]*config.CarSettings
	context  *application.Context
	defaults *config.CarSettings
}

func CreateFromConfig(ctx *application.Context) (*CarConfigRepository, error) {
	c := &config.Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	if c.Car.Defaults.MinSpeed == nil {
		percent := types.Percent(0)
		c.Car.Defaults.MinSpeed = &percent
	}

	if c.Car.Defaults.MaxBreaking == nil {
		percent := types.Percent(100)
		c.Car.Defaults.MaxBreaking = &percent
	}

	if c.Car.Defaults.PitLane == nil {
		percent := types.Percent(100)
		c.Car.Defaults.PitLane = &config.PitLaneConfig{MaxSpeed: &percent}
	}

	if c.Car.Defaults.Drivers == nil {
		c.Car.Defaults.Drivers = &types.Drivers{}
	}

	ccr := new(CarConfigRepository)
	ccr.cars = make(map[types.Id]*car.Car)
	ccr.config = make(map[types.Id]*config.CarSettings)
	ccr.defaults = c.Car.Defaults
	ccr.context = ctx

	for _, cs := range c.Car.Cars {
		ccr.config[*cs.Id] = cs
	}
	return ccr, nil
}

func (c *CarConfigRepository) Get(id types.Id, ctx ctx.Context) (*car.Car, bool, bool) {
	carCreated := false
	if _, ok := c.cars[id]; !ok {
		if _, ok := c.config[id]; !ok {
			c.config[id] = &config.CarSettings{}
		}

		if c.config[id].Drivers == nil {
			c.config[id].Drivers = &types.Drivers{
				{Name: getRandomDriver()},
			}
		}

		// i := merge.Merge(c.defaults, c.config[id])
		c.cars[id] = car.NewCar(c.context.Implement, c.context.ValueFactory, c.config[id], c.defaults, id)

		for _, r := range c.context.Rules.CarRules() {
			r.ConfigureCarState(c.cars[id], c.context.ValueFactory)
		}
		for _, r := range c.context.Rules.CarRules() {
			r.InitializeCarState(c.cars[id], ctx)
		}
		c.cars[id].Init(ctx, c.context.Postprocessors.ValuePostProcessor())
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
