package v3

import (
	"github.com/rs/zerolog"
)

const (
	RaceUnknownByte           = 0x00
	RaceStoppedByte           = 0x01
	RaceRunningByte           = 0x03
	RacePausedByte            = 0x04
	RaceFlaggedLcEnabledByte  = 0x05
	RaceFlaggedLcDisabledByte = 0x15
)

type Race struct {
	status byte
	logger zerolog.Logger
}

func NewRace(logger zerolog.Logger) *Race {
	return &Race{
		logger: logger,
		status: RaceUnknownByte,
	}
}

func (r *Race) Start() {
	r.logger.Info().
		Str("drivers", "driver3x").
		Str("race-status", "start").
		Msg("Race status changed")
	r.status = RaceRunningByte
}

func (r *Race) Flag() {
	r.status = RaceFlaggedLcEnabledByte
}

func (r *Race) Pause() {
	r.logger.Info().
		Str("drivers", "driver3x").
		Str("race-status", "pause").
		Msg("Race status changed")
	r.status = RacePausedByte
}

func (r *Race) Stop() {
	r.logger.Info().
		Str("drivers", "driver3x").
		Str("race-status", "stopped").
		Msg("Race status changed")
	r.status = RaceStoppedByte
}
