package config

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/config/context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Plugin struct {
	Name string `yaml:"plugin"`
}

func FromFile(c *context.Context, file *string) error {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		return err
	}
	config := &context.Config{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return errors.New("Failed to load config file: " + err.Error())
	}
	c.Config = config

	return nil
}

func readConfig(config []byte) (*context.Config, error) {
	c := &context.Config{}
	err := yaml.Unmarshal(config, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
