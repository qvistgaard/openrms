package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/config/application"
)

type Config struct {
	Postprocessors struct {
		InfluxDB struct {
			Url          string `mapstructure:"url"`
			BatchSize    uint   `mapstructure:"batch-size"`
			AuthToken    string `mapstructure:"auth-token"`
			Organization string `mapstructure:"organization"`
			Bucket       string `mapstructure:"bucket"`
		} `mapstructure:"influxdb"`
	} `mapstructure:"postprocessors"`
}

func CreateFromConfig(ctx *application.Context) (*Leaderboard, error) {
	return NewLeaderboard(ctx), nil
}
