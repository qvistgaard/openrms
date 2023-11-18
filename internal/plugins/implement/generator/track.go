package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
)

type Track struct {
	pitLane *PitLane
}

func NewTrack() *Track {
	return &Track{
		pitLane: NewPitLane(),
	}
}

func (t *Track) MaxSpeed(_ types.Percent) {

}

func (t *Track) PitLane() implement.PitLaneImplementer {
	return t.pitLane
}
