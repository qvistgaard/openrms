package oxigen

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Implement struct {
		Oxigen struct {
			Port string `yaml:"port"`
		} `yaml:"oxigen"`
	} `yaml:"implement"`
}

func CreateFromConfig(config []byte) (*Oxigen, error) {
	c := &Config{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}

	connection, err := CreateUSBConnection(c.Implement.Oxigen.Port)
	if err != nil {
		return nil, err
	}
	return CreateImplement(connection)
}
