package config

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/plugins/implement/oxigen"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/telemetry/influxdb"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/telemetry"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Implement Plugin `yaml:"implement"`
	Race      Plugin
	Car       Plugin
	Telemetry Plugin
	Rules     []Plugin
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
	}
	return nil, errors.New("Unknown implementer: " + c.Implement.Name)
}

func CreateTelemetryReceiverFromConfig(config []byte) (telemetry.Receiver, error) {
	c, err := readConfig(config)
	var processor telemetry.Processor
	if err != nil {
		return nil, err
	}
	var perr error
	switch c.Telemetry.Name {
	case "influxdb":
		processor, err = influxdb.CreateFromConfig(config)
	}
	if perr != nil {
		return nil, err
	}

	return telemetry.NewQueueReceiver(processor), nil
}

func CreateRaceRulesFromConfig(config []byte) ([]state.Rule, error) {
	c, err := readConfig(config)
	if err != nil {
		return nil, err
	}
	var rules []state.Rule
	for _, r := range c.Rules {
		switch r.Name {
		case "fuel":
			rules = append(rules, &fuel.Consumption{})
		}
	}
	return rules, nil
}
