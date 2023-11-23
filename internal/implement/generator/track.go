package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
)

type Track struct {
	pitLane *PitLane
}

func NewTrack() *Track {
	return &Track{
		pitLane: NewPitLane(),
	}
}

func (t *Track) MaxSpeed(percent uint8) {

}

func (t *Track) PitLane() implement.PitLaneImplementer {
	return t.pitLane
}
