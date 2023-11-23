package configuration

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config map[string]interface{}

// FromFile reads and parses a configuration file specified by the given file path.
// It takes a pointer to a file path as input and returns a `Config` map and an error, if any.
//
// The `file` parameter should point to the path of the configuration file to be read.
//
// Example usage:
//
//   filePath := "path/to/your/config.yaml"
//
//   config, err := FromFile(&filePath)
//   if err != nil {
//       log.Fatal("Failed to load configuration from file: ", err)
//   }
//
//   // Use the 'config' map for accessing configuration settings.
//
// Returns:
//   - A `Config` map representing the parsed configuration settings.
//   - An error if there was an issue reading or parsing the configuration file.
func FromFile(file *string) (Config, error) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to configuration from file")
	}
	config := Config{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, errors.New("Failed to load config file: " + err.Error())
	}

	return config, nil
}
