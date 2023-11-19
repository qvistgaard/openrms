package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type Car struct {
	id *byte
}

func NewCar(id uint8) implement.CarImplementer {
	return &Car{id: &id}
}

func (c Car) MaxSpeed(percent types.Percent) {
	log.WithField("value", percent).
		Info("Max speed updated")

}

func (c Car) PitLaneMaxSpeed(percent types.Percent) {
	log.WithField("value", percent).
		Info("Pit lane max speed updated")
}

func (c Car) MaxBreaking(percent types.Percent) {
	log.WithField("value", percent).
		Info("Maximum braking updated")

}

func (c Car) MinSpeed(percent types.Percent) {
	log.WithField("value", percent).
		Info("Min speed updated")
}
