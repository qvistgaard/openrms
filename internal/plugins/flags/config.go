package flags

type FlagConfig struct {
	Pause    *bool
	MaxSpeed *uint8
}

type Config struct {
	Plugin struct {
		Flag struct {
			Enabled bool       `mapstructure:"enabled" default:"true"`
			Yellow  FlagConfig `mapstructure:"yellow"`
			Red     FlagConfig `mapstructure:"red"`
		} `mapstructure:"flag"`
	} `mapstructure:"plugins"`
}
