package generator

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
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

func CreateFromConfig(ctx *context.Context) (*Generator, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	return &Generator{
		cars:     c.Implement.Generator.Cars,
		interval: c.Implement.Generator.Interval,
		events:   make(chan implement.Event, 1024),
	}, nil
}
