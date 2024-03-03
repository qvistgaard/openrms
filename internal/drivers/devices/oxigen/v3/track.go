package v3

import (
	"github.com/qvistgaard/openrms/internal/drivers"
)

type Track struct {
	maxSpeed uint8
	pitLane  *PitLane
}

func NewTrack() *Track {
	return &Track{
		pitLane:  NewPitLane(),
		maxSpeed: 100,
	}
}

func (t *Track) MaxSpeed(percent uint8) {
	t.maxSpeed = percentageToByte(percent)
}

func (t *Track) PitLane() drivers.PitLane {
	return t.pitLane
}
