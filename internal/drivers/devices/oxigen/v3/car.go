package v3

import (
	"fmt"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/rs/zerolog"
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

func newCar(driver *Driver3x, id types.CarId) *Car {
	return &Car{logger: driver.logger, id: id, driver3x: driver, maxBreaking: 255, minSpeed: 0, maxSpeed: 255, pitLaneSpeed: 255}
}

type Car struct {
	id           types.CarId
	driver3x     *Driver3x
	maxBreaking  uint8
	minSpeed     uint8
	maxSpeed     uint8
	pitLaneSpeed uint8
	logger       zerolog.Logger
}

func (c *Car) Id() types.CarId {
	return c.id
}

func (c *Car) SetMaxBreaking(percent uint8) {
	c.maxBreaking = percentageToByte(percent)
	c.logger.Info().
		Str("drivers", "driver3x").
		Stringer("car", c.Id()).
		Uint8("max-breaking", percent).
		Int("cmd", carMaxBreakingCode).
		Str("hex", fmt.Sprintf("%x", c.maxBreaking)).
		Msg("set car max breaking")
	c.sendMaxBreaking()
}

func (c *Car) sendMaxBreaking() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxBreakingCode, c.maxBreaking)
}

func (c *Car) SetMinSpeed(percent uint8) {
	c.minSpeed = percentageToByte(percent) >> 1
	c.logger.Info().
		Str("drivers", "driver3x").
		Stringer("car", c.Id()).
		Uint8("min-speed", percent).
		Int("cmd", carMinSpeedCode).
		Str("hex", fmt.Sprintf("%x", c.minSpeed)).
		Msg("set car min speed")
	c.sendMinSpeed()
}

func (c *Car) sendMinSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMinSpeedCode, c.minSpeed)
}

func (c *Car) SetMaxSpeed(percent uint8) {
	c.maxSpeed = percentageToByte(percent)
	c.logger.Info().
		Str("drivers", "driver3x").
		Stringer("car", c.Id()).
		Uint8("max-speed", percent).
		Int("cmd", carMaxSpeedCode).
		Str("hex", fmt.Sprintf("%x", c.maxSpeed)).
		Msg("set car max speed")
	c.sendMaxSpeed()
}

func (c *Car) sendMaxSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxSpeedCode, c.maxSpeed)
}

func (c *Car) SetPitLaneMaxSpeed(percent uint8) {
	c.pitLaneSpeed = percentageToByte(percent)
	c.logger.Info().
		Str("drivers", "driver3x").
		Stringer("car", c.Id()).
		Uint8("max-speed", percent).
		Int("cmd", carPitLaneSpeedCode).
		Str("hex", fmt.Sprintf("%x", c.pitLaneSpeed)).
		Msg("set car pit lane max speed")
	c.sendPitLaneMaxSpeed()
}

func (c *Car) sendPitLaneMaxSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carPitLaneSpeedCode, c.pitLaneSpeed)
}
