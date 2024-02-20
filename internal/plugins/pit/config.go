package pit

type Config struct {
	Plugin struct {
		Pit struct {
			Enabled    bool `mapstructure:"enabled"`
			Commentary bool `mapstructure:"announcements" default:"true"`
		} `mapstructure:"pit"`
	} `mapstructure:"plugins"`
}
