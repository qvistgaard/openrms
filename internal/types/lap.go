package types

import (
	"time"
)

type Lap struct {
	Number    uint32
	Time      time.Duration
	RaceTimer time.Duration
}
