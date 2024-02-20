package ontrack

type Config struct {
	Plugin struct {
		OnTrack struct {
			Enabled    bool   `mapstructure:"enabled"`
			Flag       string `mapstructure:"flag" default:"red"`
			Commentary bool   `mapstructure:"announcements" default:"true"`
		} `mapstructure:"ontrack"`
	} `mapstructure:"plugins"`
}
