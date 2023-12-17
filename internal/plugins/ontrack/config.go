package ontrack

type Config struct {
	Plugin struct {
		OnTrack struct {
			Enabled bool   `mapstructure:"enabled"`
			Flag    string `mapstructure:"flag" default:"red"`
		} `mapstructure:"ontrack"`
	} `mapstructure:"plugins"`
}
