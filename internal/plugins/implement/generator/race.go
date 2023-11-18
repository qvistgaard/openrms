package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
)

type Race struct {
}

func NewRace() *Race {
	return &Race{}
}

func (r *Race) Status(status implement.RaceStatus) {

}
