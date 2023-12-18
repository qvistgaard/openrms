package track

type Config struct {
	Track struct {
		MaxSpeed uint8 `mapstructure:"max-speed"`
		PitLane  struct {
			LapCounting struct {
				Enabled bool
				OnEntry bool `mapstructure:"on-entry"`
			} `mapstructure:"lap-counting"`
		} `mapstructure:"pit-lane"`
	}
}
