package pit

type Config struct {
	Plugin struct {
		Pit struct {
			Enabled bool `mapstructure:"enabled"`
		} `mapstructure:"pit"`
	} `mapstructure:"plugins"`
}
