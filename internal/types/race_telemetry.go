package types

import (
	"sort"
	"time"
)

type RaceTelemetryEntry struct {
	Car         CarId
	Laps        Lap
	Delta       time.Duration
	Best        time.Duration
	Last        time.Duration
	Deslotted   bool
	MinSpeed    uint8
	MaxSpeed    uint8
	MaxPitSpeed uint8
	InPit       bool
	LimbMode    bool
	Fuel        float32
	Name        string
}

type RaceTelemetry map[CarId]*RaceTelemetryEntry

func (r RaceTelemetry) Sort() []RaceTelemetryEntry {
	sorted := make([]RaceTelemetryEntry, 0, len(r))
	for _, v := range r {
		sorted = append(sorted, *v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Laps.Number > sorted[j].Laps.Number {
			return true
		} else if sorted[i].Laps.Number == sorted[j].Laps.Number {
			if sorted[i].Laps.RaceTimer < sorted[j].Laps.RaceTimer {
				return true
			}
			return false
		} else {
			return false
		}
	})
	return sorted
}
