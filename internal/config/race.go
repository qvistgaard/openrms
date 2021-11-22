package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/track"
	"github.com/qvistgaard/openrms/internal/types"
)

type TrackConfig struct {
	Track struct {
		MaxSpeed types.Percent `mapstructure:"max-speed"`
		PitLane  struct {
			LapCounting struct {
				Enabled bool
				OnEntry bool `mapstructure:"on-entry"`
			} `mapstructure:"lap-counting"`
		} `mapstructure:"pit-lane"`
	}
}

func ConfigureTrack(ctx *application.Context) error {
	c := &TrackConfig{}

	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil
	}

	var o implement.PitLaneLapCounting
	if c.Track.PitLane.LapCounting.OnEntry {
		o = implement.LapCountingOnEntry
	} else {
		o = implement.LapCountingOnExit
	}

	ctx.Implement.Track().PitLane().LapCounting(c.Track.PitLane.LapCounting.Enabled, o)
	ctx.Implement.Track().MaxSpeed(c.Track.MaxSpeed)
	ctx.Track = track.NewTrack(ctx.Implement)

	return nil
}
