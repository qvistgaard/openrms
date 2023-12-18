package generator

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

func NewCar(id types.CarId, lap uint32) drivers.Car {
	return &Car{id, lap}
}

type Car struct {
	id  types.CarId
	lap uint32
}

func (c Car) Id() types.CarId {
	return c.id
}

func (c Car) SetMaxSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", c.id).
		Info("Max speed updated")

}

func (c Car) SetPitLaneMaxSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", c.id).
		Info("Pit lane max speed updated")
}

func (c Car) SetMaxBreaking(percent uint8) {
	log.WithField("value", percent).
		WithField("car", c.id).
		Info("Maximum braking updated")

}

func (c Car) SetMinSpeed(percent uint8) {
	log.WithField("value", percent).
		WithField("car", c.id).
		Info("Min speed updated")
}
