package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
)

type Config struct {
	Implement struct {
		Generator struct {
			Cars     uint8
			Interval uint
		}
	}
}

func New(c Config) (implement.Implementer, error) {
	return &Generator{
		cars:     c.Implement.Generator.Cars,
		interval: c.Implement.Generator.Interval,
		events:   make(chan implement.Event, 1024),
		race:     NewRace(),
	}, nil
}
