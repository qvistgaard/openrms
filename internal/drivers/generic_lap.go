package drivers

import (
	"time"
)

func GenericLap(lapNumber uint16, lapTime time.Duration, recordTime time.Duration) Lap {
	return genericLap{
		lapNumber:  lapNumber,
		lapTime:    lapTime,
		recordTime: recordTime,
	}
}

type genericLap struct {
	lapNumber  uint16
	lapTime    time.Duration
	recordTime time.Duration
}

func (g genericLap) Number() uint16 {
	return g.lapNumber
}

func (g genericLap) Time() time.Duration {
	return g.lapTime
}

func (g genericLap) Recorded() time.Duration {
	return g.recordTime
}
