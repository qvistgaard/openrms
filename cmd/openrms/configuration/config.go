package configuration

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

//go:embed files/config.yaml
var files embed.FS

type Config map[string]interface{}

// FromFile reads and parses a configuration file specified by the given file path.
// It takes a pointer to a file path as input and returns a `Config` map and an error, if any.
//
// The `file` parameter should point to the path of the configuration file to be read.
//
// Example usage:
//
//	filePath := "path/to/your/config.yaml"
//
//	config, err := FromFile(&filePath)
//	if err != nil {
//	    log.Fatal("Failed to load configuration from file: ", err)
//	}
//
//	// Use the 'config' map for accessing configuration settings.
//
// Returns:
//   - A `Config` map representing the parsed configuration settings.
//   - An error if there was an issue reading or parsing the configuration file.
func FromFile(file *string) (Config, error) {
	err := createConfigFileIfNotExists(file)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(*file)
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

func createConfigFileIfNotExists(config *string) error {
	if _, err := os.Stat(*config); errors.Is(err, os.ErrNotExist) {
		answer := "!"
		input := bufio.NewScanner(os.Stdin)

		for answer != "y" && answer != "n" && answer != "" {
			fmt.Println("Config file not found. Would you like to create it? [Y/n]")
			input.Scan()
			answer = strings.ToLower(input.Text())
			if answer == "y" || answer == "" {

				file, err := files.ReadFile("files/config.yaml")
				if err != nil {
					return err
				}
				err = os.WriteFile(*config, file, 0666)
				if err != nil {
					return err
				}
				log.Info("Config file created")
			}
		}
	}
	return nil

}
