package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Implement struct {
		Generator struct {
			Cars     uint8
			Interval uint
		}
	}
}

func CreateFromConfig(config []byte) (*Generator, error) {
	c := &Config{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}

	return &Generator{
		cars:     c.Implement.Generator.Cars,
		interval: c.Implement.Generator.Interval,
		events:   make(chan implement.Event, 1024),
	}, nil
}
