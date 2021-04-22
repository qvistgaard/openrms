package config

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/implement"
	carConfig "github.com/qvistgaard/openrms/internal/plugins/car/config"
	"github.com/qvistgaard/openrms/internal/plugins/implement/generator"
	"github.com/qvistgaard/openrms/internal/plugins/implement/oxigen"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/influxdb"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/websocket"
	"github.com/qvistgaard/openrms/internal/plugins/rules/damage"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/rules/leaderboard"
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/plugins/rules/tirewear"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Implement      Plugin
	Race           Plugin
	Car            Plugin
	Telemetry      Plugin
	Rules          []Plugin
	PostProcessors map[string]interface{}
}

type Plugin struct {
	Name string `yaml:"plugin"`
}

func readConfig(config []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(config, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func CreateImplementFromConfig(config []byte) (implement.Implementer, error) {
	c, err := readConfig(config)
	if err != nil {
		return nil, err
	}

	switch c.Implement.Name {
	case "oxigen":
		return oxigen.CreateFromConfig(config)
	case "generator":
		return generator.CreateFromConfig(config)
	}
	return nil, errors.New("Unknown implementer: " + c.Implement.Name)
}

func CreatePostProcessors(config []byte) ([]postprocess.PostProcessor, error) {
	c, err := readConfig(config)
	if err != nil {
		return nil, err
	}
	var p []postprocess.PostProcessor
	for k, _ := range c.PostProcessors {
		switch k {
		case "influxdb":
			influxdb, err := influxdb.CreateFromConfig(config)
			if err != nil {
				return nil, err
			}
			p = append(p, influxdb)
			go influxdb.Process()
		case "websocket":
			websocket, err := websocket.CreateFromConfig(config)
			if err != nil {
				return nil, err
			}
			p = append(p, websocket)
			go websocket.Process()
		}
	}
	return p, nil
}

func CreateRaceRulesFromConfig(config []byte) (state.Rules, error) {
	c, err := readConfig(config)
	if err != nil {
		return nil, err
	}
	rules := state.CreateRuleList()
	for _, r := range c.Rules {
		switch r.Name {
		case "fuel":
			rules.Append(&fuel.Consumption{})
		case "limb-mode":
			rules.Append(&limbmode.LimbMode{})
		case "damage":
			rules.Append(&damage.Damage{})
		case "pit":
			rules.Append(pit.CreatePitRule(rules))
		case "tirewear":
			rules.Append(&tirewear.TireWear{})
		case "leaderboard":
			rules.Append(&leaderboard.BoardDefault{})
		default:
			return nil, errors.New("Unknown rule: " + r.Name)
		}
	}
	return rules, nil
}

func CreateCarRepositoryFromConfig(config []byte) (car.Repository, error) {
	c, err := readConfig(config)
	if err != nil {
		return nil, err
	}

	switch c.Car.Name {
	case "config":
		return carConfig.CreateFromConfig(config)
	}
	return nil, errors.New("no car configuration found")
}

func CreateCourseFromConfig(config []byte, rules state.Rules) (*state.Course, error) {
	return state.CreateCourseFromConfig(config, rules)
}
