package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
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
	t.maxSpeed = percent
}

func (t *Track) PitLane() implement.PitLaneImplementer {
	return t.pitLane
}
