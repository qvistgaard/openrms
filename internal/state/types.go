package state

import (
	"math"
	"time"
)

type Percent float64

func PercentFromUint8(i uint8) Percent {
	return Percent(math.Round(float64(i) / 255 * 100))
}

func PercentToUint8(f float64) uint8 {
	return uint8(255 * (Percent(f) / 100))
}

type CarId uint8
type Speed Percent
type TriggerValue Percent
type Breaking Percent
type LapNumber uint
type RaceTimer time.Duration
type LapTime time.Duration
