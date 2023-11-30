package generator

import (
	"github.com/qvistgaard/openrms/internal/drivers"
)

type Config struct {
	Implement struct {
		Generator struct {
			Cars     uint8
			Interval uint
		}
	}
}

func New(c Config) (drivers.Driver, error) {
	return &Generator{
		cars:     c.Implement.Generator.Cars,
		interval: c.Implement.Generator.Interval,
		events:   make(chan drivers.Event, 1024),
		race:     NewRace(),
	}, nil
}
