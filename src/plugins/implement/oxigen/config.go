package oxigen

import (
	"gopkg.in/yaml.v2"
)

// TODO: Config structure is incorrect
type Config struct {
	port string
}

func CreateImplementFromConfig(config string) (*Oxigen, error) {
	c := &Config{}
	perr := yaml.Unmarshal([]byte(config), c)
	if perr != nil {
		return nil, perr
	}

	connection, err := CreateUSBConnection(c.port)
	if err != nil {
		return nil, err
	}
	return CreateImplement(connection)
}
