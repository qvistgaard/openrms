package v3

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

func newCar(driver *Driver3x, id types.CarId) drivers.Car {
	return Car{id, driver}
}

type Car struct {
	id       types.CarId
	driver3x *Driver3x
}

func (c Car) Id() types.CarId {
	return c.id
}

func (c Car) SetMaxBreaking(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-breaking", percent).
		WithField("cmd", carMaxBreakingCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max breaking")
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxBreakingCode, toByte)
}

func (c Car) SetMinSpeed(percent uint8) {
	toByte := percentageToByte(percent) >> 1
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("min-speed", percent).
		WithField("cmd", carMinSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car min speed")
	c.driver3x.sendCarCommand(uint8(c.Id()), carMinSpeedCode, toByte)
}

func (c Car) SetMaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carMaxSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car max speed")
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxSpeedCode, toByte)
}

func (c Car) SetPitLaneMaxSpeed(percent uint8) {
	toByte := percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carPitLaneSpeedCode).
		WithField("hex", fmt.Sprintf("%x", toByte)).
		Info("set car pit lane max speed")
	c.driver3x.sendCarCommand(uint8(c.Id()), carPitLaneSpeedCode, toByte)
}
