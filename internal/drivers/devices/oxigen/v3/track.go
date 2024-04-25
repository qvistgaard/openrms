package v3

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/rs/zerolog"
)

type Track struct {
	maxSpeed uint8
	pitLane  *PitLane
	logger   zerolog.Logger
}

func NewTrack(logger zerolog.Logger) *Track {
	return &Track{
		logger:   logger,
		pitLane:  NewPitLane(logger),
		maxSpeed: 100,
	}
}

func (t *Track) MaxSpeed(percent uint8) {
	t.logger.Info().
		Str("drivers", "driver3x").
		Uint8("max-speed", percent).
		Msg("Track max speed changed")
	t.maxSpeed = percentageToByte(percent)
}

func (t *Track) PitLane() drivers.PitLane {
	return t.pitLane
}
