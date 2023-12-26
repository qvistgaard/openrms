package telemetry

import (
	"github.com/qvistgaard/openrms/internal/types"
	"sort"
	"time"
)

type Entry struct {
	Id              uint
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
	Enabled         bool
	PitStopSequence uint8
	Color           string
}

type Race map[types.CarId]*Entry

func (r Race) Sort() []Entry {
	sorted := make([]Entry, 0, len(r))
	for _, v := range r {
		sorted = append(sorted, *v)
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if !sorted[i].Enabled && sorted[j].Enabled {
			return false
		} else if sorted[i].Enabled && !sorted[j].Enabled {
			return true
		}

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
