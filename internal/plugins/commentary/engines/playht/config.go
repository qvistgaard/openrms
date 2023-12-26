package playht

type PlayHTConfig struct {
	Voice  string  `mapstructure:"voice" default:"Oliver (Advertising)"`
	ApiKey string  `mapstructure:"apiKey"`
	UserId string  `mapstructure:"userId"`
	Speed  float32 `mapstructure:"speed" default:"1.1"`
	Cache  string  `mapstructure:"cache" default:"cache"`
}
