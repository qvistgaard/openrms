package pit

type Config struct {
	Plugin struct {
		Pit struct {
			Enabled    bool `mapstructure:"enabled"`
			Commentary bool `mapstructure:"commentary" default:"true"`
		} `mapstructure:"pit"`
	} `mapstructure:"plugins"`
}
