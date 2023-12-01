package generator

import (
	"github.com/qvistgaard/openrms/internal/drivers"
)

type PitLane struct {
}

func NewPitLane() *PitLane {
	return &PitLane{}
}

func (p *PitLane) LapCounting(_ bool, _ drivers.PitLaneLapCounting) {

}
