package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
)

// Leaving for later when implementing lane change limitations.
const (
	CarForceLaneChangeLeft  = 0x80
	CarForceLaneChangeRight = 0x40
	CarForceLaneChangeNone  = 0x00
	CarForceLangeChangeAny  = CarForceLaneChangeLeft | CarForceLaneChangeRight
)

const (
	carMaxSpeedCode     = 0x82
	carPitLaneSpeedCode = 0x81
	carMinSpeedCode     = 0x03
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
	c.oxigen.sendCarCommand(c.id, carMaxBreakingCode, percent.Uint8())
}

func (c *Car) MinSpeed(percent types.Percent) {
	c.oxigen.sendCarCommand(c.id, carMinSpeedCode, percent.Uint8())
}

func (c *Car) MaxSpeed(percent types.Percent) {
	c.oxigen.sendCarCommand(c.id, carMaxSpeedCode, percent.Uint8())
}

func (c *Car) PitLaneMaxSpeed(percent types.Percent) {
	c.oxigen.sendCarCommand(c.id, carPitLaneSpeedCode, percent.Uint8())
}
