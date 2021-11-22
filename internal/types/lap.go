package types

import (
	"github.com/qvistgaard/openrms/internal/telemetry"
	"time"
)

type Lap struct {
	LapNumber uint16
	RaceTimer time.Duration
}

func NewLap(lapNumber uint16, raceTimer time.Duration) Lap {
	return Lap{LapNumber: lapNumber, RaceTimer: raceTimer}
}

func (l Lap) Metrics() []telemetry.Metric {
	return []telemetry.Metric{
		{
			Name:  "lap-number",
			Value: l.LapNumber,
		},
	}
}
