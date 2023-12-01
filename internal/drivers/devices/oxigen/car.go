package oxigen

import (
	"fmt"
	"github.com/qvistgaard/openrms/internal/drivers"
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

func NewCar(oxigen *Oxigen, id types.Id) drivers.Car {
	return Car{id, oxigen}
}

type Car struct {
	id     types.Id
	oxigen *Oxigen
}

func (c Car) Id() types.Id {
	return c.id
}

func (c Car) SetMaxBreaking(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "oxigen").
		WithField("car", c.Id()).
		WithField("max-breaking", percent).
		WithField("cmd", carMaxBreakingCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max breaking")
	c.oxigen.sendCarCommand(uint8(c.Id()), carMaxBreakingCode, toByte)
}

func (c Car) SetMinSpeed(percent uint8) {
	toByte := percentageToByte(percent) >> 1
	log.WithField("drivers", "oxigen").
		WithField("car", c.Id()).
		WithField("min-speed", percent).
		WithField("cmd", carMinSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car min speed")
	c.oxigen.sendCarCommand(uint8(c.Id()), carMinSpeedCode, toByte)
}

func (c Car) SetMaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "oxigen").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carMaxSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max speed")
	c.oxigen.sendCarCommand(uint8(c.Id()), carMaxSpeedCode, toByte)
}

func (c Car) SetPitLaneMaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "oxigen").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carPitLaneSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car pit lane max speed")
	c.oxigen.sendCarCommand(uint8(c.Id()), carPitLaneSpeedCode, toByte)
}

func percentageToByte(percent uint8) uint8 {
	if percent > 100 {
		percent = 100
	}
	return uint8(255.0 * (float64(percent) / 100.0))
}
