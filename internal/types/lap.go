package types

import (
	"time"
)

type Lap struct {
	Number    uint16
	Time      time.Duration
	RaceTimer time.Duration
}
