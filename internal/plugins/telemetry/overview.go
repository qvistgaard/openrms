package telemetry

import (
	"github.com/qvistgaard/openrms/internal/types"
	"sort"
	"time"
)

type Entry struct {
	Id              types.CarId
	Team            string
	Last            types.Lap
	PitStopActive   bool
	Laps            []types.Lap
	Delta           time.Duration
	Best            time.Duration
	Deslotted       bool
	OnTrack         bool
	MinSpeed        uint8
	MaxSpeed        uint8
	MaxPitSpeed     uint8
	InPit           bool
	LimbMode        bool
	Fuel            float32
	PitStopSequence uint8
}

type Race map[types.CarId]*Entry

func (r Race) Sort() []Entry {
	sorted := make([]Entry, 0, len(r))
	for _, v := range r {
		sorted = append(sorted, *v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Last.Number > sorted[j].Last.Number {
			return true
		} else if sorted[i].Last.Number == sorted[j].Last.Number {
			if sorted[i].Last.RaceTimer < sorted[j].Last.RaceTimer {
				return true
			}
			return false
		} else {
			return false
		}
	})
	return sorted
}