package types

import (
	"sort"
	"time"
)

type RaceTelemetryEntry struct {
	Car         Id
	Laps        Lap
	Delta       time.Duration
	Best        time.Duration
	Last        time.Duration
	Deslotted   bool
	MinSpeed    float64
	MaxSpeed    float64
	MaxPitSpeed float64
	InPit       bool
	Fuel        float64
	Name        string
	PitState    CarPitState
}

type RaceTelemetry map[Id]*RaceTelemetryEntry

func (r RaceTelemetry) Sort() []RaceTelemetryEntry {
	sorted := make([]RaceTelemetryEntry, 0, len(r))
	for _, v := range r {
		sorted = append(sorted, *v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Laps.LapNumber > sorted[j].Laps.LapNumber {
			return true
		} else if sorted[i].Laps.LapNumber == sorted[j].Laps.LapNumber {
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
