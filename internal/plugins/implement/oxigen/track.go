package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
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

func (t *Track) MaxSpeed(maxSpeed types.Percent) {
	t.maxSpeed = maxSpeed.Uint8()
}

func (t *Track) PitLane() implement.PitLaneImplementer {
	return t.pitLane
}
