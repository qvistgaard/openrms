package elevenlabs

type ElevenLabsConfig struct {
	Voice  string  `mapstructure:"voice" default:"Patrick"`
	ApiKey string  `mapstructure:"apiKey"`
	Speed  float32 `mapstructure:"speed" default:"1.1"`
	Cache  string  `mapstructure:"cache" default:"cache"`
}
