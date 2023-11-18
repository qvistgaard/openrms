package web

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/types"
)

type Config struct {
	Car struct {
		Cars []struct {
			Id      types.Id
			Drivers []struct {
				Name string
			}
		}
	}
}

func CreateFromConfig(ctx *application.Context) (*Leaderboard, error) {
	c := &Config{}
	mapstructure.Decode(ctx.Config, c)

	return NewLeaderboard(ctx, c), nil
}
