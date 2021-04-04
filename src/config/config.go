package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"openrms/implement"
	"openrms/plugins/implement/oxigen"
)

type Config struct {
	implement struct {
		plugin string
	}
	race struct {
		plugin string
	}
	cars struct {
		plugin string
	}
	telemetry struct {
		plugin string
	}
}

//TODO: Do this smarter only reading config once.
func CreateImplementFromConfig(config string) (implement.Implementer, error) {
	c := &Config{}
	yaml.Unmarshal([]byte(config), c)

	switch c.implement.plugin {
	case "oxigen":
		return oxigen.CreateImplementFromConfig(config)
	}
	return nil, errors.New("Unknown implementer: " + c.implement.plugin)

}
