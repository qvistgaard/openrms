package config

import (
	"errors"
	"github.com/qvistgaard/openrms/internal/config/application"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Plugin struct {
	Name string `yaml:"plugin"`
}

func FromFile(c *application.Context, file *string) error {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		return err
	}
	config := &application.Config{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return errors.New("Failed to load config file: " + err.Error())
	}
	c.Config = config

	return nil
}
