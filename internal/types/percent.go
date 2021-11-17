package types

import "math"

type Percent float64

func NewPercentFromUint8(i uint8) Percent {
	return NewPercentFromFloat64(float64(i))
}

func NewPercentFromFloat64(i float64) Percent {
	return Percent(math.Round(i / 255 * 100))
}

func (p Percent) Uint8() uint8 {
	return uint8(255 * (p / 100))
}
