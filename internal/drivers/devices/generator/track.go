package generator

import (
	"github.com/qvistgaard/openrms/internal/drivers"
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

func (t *Track) PitLane() drivers.PitLane {
	return t.pitLane
}
