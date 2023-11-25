package oxigen

import (
	"fmt"
	"github.com/qvistgaard/openrms/internal/implement"
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

func (c *Car) MaxBreaking(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-breaking", percent).
		WithField("cmd", carMaxBreakingCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max breaking")
	c.oxigen.sendCarCommand(c.id, carMaxBreakingCode, toByte)
}

func (c *Car) MinSpeed(percent uint8) {
	toByte := percentageToByte(percent) >> 1
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("min-speed", percent).
		WithField("cmd", carMinSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car min speed")
	c.oxigen.sendCarCommand(c.id, carMinSpeedCode, toByte)
}

func (c *Car) MaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-speed", percent).
		WithField("cmd", carMaxSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max speed")
	c.oxigen.sendCarCommand(c.id, carMaxSpeedCode, toByte)
}

func (c *Car) PitLaneMaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("implement", "oxigen").
		WithField("car", *c.id).
		WithField("max-speed", percent).
		WithField("cmd", carPitLaneSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car pit lane max speed")
	c.oxigen.sendCarCommand(c.id, carPitLaneSpeedCode, toByte)
}

func percentageToByte(percent uint8) uint8 {
	if percent > 100 {
		percent = 100
	}
	return uint8(255.0 * (float64(percent) / 100.0))
}
