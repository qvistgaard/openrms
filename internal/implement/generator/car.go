package generator

import (
	"github.com/qvistgaard/openrms/internal/implement"
	log "github.com/sirupsen/logrus"
)

type Car struct {
	id *byte
}

func NewCar(id uint8) implement.CarImplementer {
	return &Car{id: &id}
}

func (c Car) MaxSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", *c.id).
		Info("Max speed updated")

}

func (c Car) PitLaneMaxSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", *c.id).
		Info("Pit lane max speed updated")
}

func (c Car) MaxBreaking(percent uint8) {
	log.WithField("value", percent).
		WithField("car", *c.id).
		Info("Maximum braking updated")

}

func (c Car) MinSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", *c.id).
		Info("Min speed updated")
}
