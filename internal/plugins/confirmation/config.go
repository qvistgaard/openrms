package confirmation

import "time"

type Config struct {
	Plugin struct {
		Confirmation struct {
			Enabled bool           `mapstructure:"enabled"`
			Timeout *time.Duration `mapstructure:"timeout"`
			Modes   []string       `mapstructure:"modes"`
		} `mapstructure:"confirmation"`
	} `mapstructure:"plugins"`
}
