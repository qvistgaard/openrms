package types

import "math"

type Percent float64

func NewPercentFromUint8(i uint8) Percent {
	return Percent(math.Round(float64(i) / 255 * 100))
}

func (p Percent) Uint8() uint8 {
	return uint8(255 * (p / 100))
}
