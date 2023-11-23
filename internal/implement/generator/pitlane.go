package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
)

type PitLane struct {
}

func NewPitLane() *PitLane {
	return &PitLane{}
}

func (p *PitLane) LapCounting(enabled bool, option implement.PitLaneLapCounting) {

}
