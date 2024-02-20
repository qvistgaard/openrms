package confirmation

import "time"

type Config struct {
	Plugin struct {
		Confirmation struct {
			Enabled       bool           `mapstructure:"enabled"`
			Timeout       *time.Duration `mapstructure:"timeout"`
			Modes         []string       `mapstructure:"modes"`
			Announcements bool           `mapstructure:"announcements"`
		} `mapstructure:"confirmation"`
	} `mapstructure:"plugins"`
}
