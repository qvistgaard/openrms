package generator

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func NewCar(id uint8, lap uint16) drivers.Car {
	return &Car{id, lap}
}

type Car struct {
	id  byte
	lap uint16
}

func (c Car) Controller() drivers.Controller {
	return controller{}
}

func (c Car) Id() types.Id {
	return types.IdFromUint(c.id)
}

func (c Car) Reset() bool {
	return false
}

func (c Car) InPit() bool {
	return false
}

func (c Car) Deslotted() bool {
	return false
}

func (c Car) Lap() drivers.Lap {
	lt := time.Duration(rand.Intn(10000)) * time.Millisecond
	return drivers.GenericLap(c.lap, lt, 0)
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
