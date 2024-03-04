package v3

import (
	"fmt"
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

func newCar(driver *Driver3x, id types.CarId) *Car {
	return &Car{id: id, driver3x: driver, maxBreaking: 255, maxSpeed: 255 >> 1, pitLaneSpeed: 255}
}

type Car struct {
	id           types.CarId
	driver3x     *Driver3x
	maxBreaking  uint8
	minSpeed     uint8
	maxSpeed     uint8
	pitLaneSpeed uint8
}

func (c *Car) Id() types.CarId {
	return c.id
}

func (c *Car) SetMaxBreaking(percent uint8) {
	c.maxBreaking = percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-breaking", percent).
		WithField("cmd", carMaxBreakingCode).
		WithField("hex", fmt.Sprintf("%x", c.maxBreaking)).
		Info("set car max breaking")
	c.sendMaxBreaking()
}

func (c *Car) sendMaxBreaking() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxBreakingCode, c.maxBreaking)
}

func (c *Car) SetMinSpeed(percent uint8) {
	c.minSpeed = percentageToByte(percent) >> 1
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("min-speed", percent).
		WithField("cmd", carMinSpeedCode).
		WithField("hex", fmt.Sprintf("%x", c.minSpeed)).
		Info("set car min speed")
	c.sendMinSpeed()
}

func (c *Car) sendMinSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMinSpeedCode, c.minSpeed)
}

func (c *Car) SetMaxSpeed(percent uint8) {
	c.maxSpeed = percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carMaxSpeedCode).
		WithField("hex", fmt.Sprintf("%x", c.maxSpeed)).
		Info("set car max speed")
	c.sendMaxSpeed()
}

func (c *Car) sendMaxSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carMaxSpeedCode, c.maxSpeed)
}

func (c *Car) SetPitLaneMaxSpeed(percent uint8) {
	c.pitLaneSpeed = percentageToByte(percent)
	log.WithField("drivers", "driver3x").
		WithField("car", c.Id()).
		WithField("max-speed", percent).
		WithField("cmd", carPitLaneSpeedCode).
		WithField("hex", fmt.Sprintf("%x", c.pitLaneSpeed)).
		Info("set car pit lane max speed")
	c.sendPitLaneMaxSpeed()
}

func (c *Car) sendPitLaneMaxSpeed() {
	c.driver3x.sendCarCommand(uint8(c.Id()), carPitLaneSpeedCode, c.pitLaneSpeed)
}
