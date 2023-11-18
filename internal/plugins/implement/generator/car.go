package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
)

type Car struct {
	id *byte
}

func NewCar(id uint8) implement.CarImplementer {
	return &Car{id: &id}
}

func (c Car) MaxSpeed(percent types.Percent) {

}

func (c Car) PitLaneMaxSpeed(percent types.Percent) {

}

func (c Car) MaxBreaking(percent types.Percent) {

}

func (c Car) MinSpeed(percent types.Percent) {

}
