package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state/race"
)

type RaceConfig struct {
}

func ConfigureRace(ctx *application.Context) error {
	c := &RaceConfig{}

	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil
	}

	ctx.Race = race.NewRace(ctx.Implement, ctx.ValueFactory)

	return nil
}
