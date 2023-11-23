package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

// Leaving for later when implementing lane change limitations.
const (
	CarForceLaneChangeLeft  = 0x80
	CarForceLaneChangeRight = 0x40
	CarForceLaneChangeNone  = 0x00
	CarForceLangeChangeAny  = CarForceLaneChangeLeft | CarForceLaneChangeRight
)

const (
	carMaxSpeedCode     = 0x02 // 0x82
	carPitLaneSpeedCode = 0x01 // 0x81
	carMinSpeedCode     = 0x03 //
	carMaxBreakingCode  = 0x05
)

func NewCar(implement *Oxigen, id uint8) implement.CarImplementer {
	return &Car{id: &id, oxigen: implement}
}

type Car struct {
	id     *byte
	oxigen *Oxigen
}

func (c *Car) MaxBreaking(percent types.Percent) {
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-breaking", percent).
		Debug("set car max breaking")
	c.oxigen.sendCarCommand(c.id, carMaxBreakingCode, percent.Uint8())
}

func (c *Car) MinSpeed(percent uint8) {
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("min-speed", percent).
		Debug("set car min speed")
	c.oxigen.sendCarCommand(c.id, carMinSpeedCode, percent>>1)
}

func (c *Car) MaxSpeed(percent uint8) {
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-speed", percent).
		WithField("max-speed-uint", percent).
		Debug("set car max speed")
	c.oxigen.sendCarCommand(c.id, carMaxSpeedCode, percent)
}

func (c *Car) PitLaneMaxSpeed(percent uint8) {
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-speed", percent).
		Debug("set car pit lane max speed")
	c.oxigen.sendCarCommand(c.id, carPitLaneSpeedCode, percent)
}
