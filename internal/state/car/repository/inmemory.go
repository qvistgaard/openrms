package repository

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/plugins"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/car/names"
	"github.com/qvistgaard/openrms/internal/types"
)

type InMemory struct {
	cars      map[types.CarId]*car.Car
	config    map[types.CarId]*car.Settings
	defaults  *car.Settings
	plugins   plugins.List
	implement drivers.Driver
}

func New(config car.Config, driver drivers.Driver, plugins plugins.List) Repository {
	if config.Car.Defaults.MinSpeed == nil {
		percent := uint8(0)
		config.Car.Defaults.MinSpeed = &percent
	}

	if config.Car.Defaults.MaxBreaking == nil {
		percent := uint8(100)
		config.Car.Defaults.MaxBreaking = &percent
	}

	if config.Car.Defaults.PitLane == nil {
		percent := uint8(100)
		config.Car.Defaults.PitLane = &car.PitLaneConfig{MaxSpeed: &percent}
	}

	if config.Car.Defaults.Drivers == nil {
		config.Car.Defaults.Drivers = &types.Drivers{}
	}

	ccr := new(InMemory)
	ccr.cars = make(map[types.CarId]*car.Car)
	ccr.config = make(map[types.CarId]*car.Settings)
	ccr.defaults = config.Car.Defaults
	ccr.plugins = plugins
	ccr.implement = driver

	for _, cs := range config.Car.Cars {
		ccr.config[*cs.Id] = cs
	}
	return ccr
}

func (c *InMemory) Get(id types.CarId) (*car.Car, bool, bool) {
	carCreated := false
	if _, ok := c.cars[id]; !ok {
		if _, ok := c.config[id]; !ok {
			c.config[id] = &car.Settings{}
		}

		if c.config[id].Drivers == nil {
			c.config[id].Drivers = &types.Drivers{
				{Name: names.RandomDriver(id)},
			}
		}
		if c.config[id].Team == nil {
			team := names.RandomTeam(id)
			c.config[id].Team = &team
		}

		if c.config[id].Manufacturer == nil {
			manufacturer := names.RandomManufacture(id)
			c.config[id].Manufacturer = &manufacturer
		}
		if c.config[id].Color == nil {
			color := ""
			c.config[id].Color = &color
		}

		// i := merge.Merge(c.defaults, c.config[id])
		c.cars[id] = car.NewCar(c.implement, c.config[id], c.defaults, id)

		for _, r := range c.plugins.Car() {
			r.ConfigureCar(c.cars[id])
		}

		c.cars[id].Initialize()
		for _, r := range c.plugins.Car() {
			r.InitializeCar(c.cars[id])
		}
		carCreated = true
	}
	return c.cars[id], true, carCreated
}
