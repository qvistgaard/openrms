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

func NewCar(implement *Oxigen, b []byte) drivers.Car {
	return Car{b, Controller(b), implement}
}

type Car struct {
	data       []byte
	controller drivers.Controller
	oxigen     *Oxigen
}

func (c Car) Id() types.Id {
	return types.IdFromUint(c.data[1])
}

func (c Car) Reset() bool {
	return 0x01&c.data[0] == 0x01
}

func (c Car) InPit() bool {
	return 0x40&c.data[8] == 0x40
}

func (c Car) Deslotted() bool {
	return !(0x80&c.data[7] == 0x80)
}

func (c Car) Controller() drivers.Controller {
	return c.controller
}

func (c Car) Lap() drivers.Lap {
	//TODO implement me
	panic("implement me")
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
