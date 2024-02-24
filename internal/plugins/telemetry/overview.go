package telemetry

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types"
	"sort"
	"time"
)

type Entry struct {
	Id              types.CarId
	Team            string
	Last            types.Lap
	Number          uint
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
	DeslotsTotal    uint32
	DeslotsLap      uint32
	car             *car.Car
}

type Race map[types.CarId]*Entry

// Sort organizes the race entries based on their position and timing, prioritizing enabled entries.
//
// The sorting criteria are as follows:
// 1. Entries that are enabled (`Enabled == true`) are prioritized over those that are not enabled.
// 2. Among enabled entries, those with a greater 'Last.Number' (representing the position or lap number) are prioritized.
// 3. If 'Last.Number' is the same between two entries, the one with a smaller 'Last.RaceTimer'
//    (indicating a faster time for that lap/position) is prioritized.
//
// Returns:
// - A slice of `Entry` structs sorted according to the criteria above.
//
// This function is useful for determining the current standings of the race, especially when
// entries need to be displayed in order of their race position and timing.

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

// FastestLap identifies the race entries with the fastest single lap times, prioritizing enabled entries.
//
// The sorting criteria are as follows:
// 1. Entries that are enabled (`Enabled == true`) are prioritized over those that are not enabled.
// 2. Among enabled entries, those with a smaller 'Best' value (representing the fastest lap time) are prioritized.
//
// Returns:
// - A slice of `Entry` structs sorted by the fastest lap times, with enabled entries taking precedence.
//
// This function is particularly useful for highlighting performances within the race, such as
// awarding a fastest lap bonus or for statistical analysis post-race. It allows for a quick
// identification of which entries had the best single-lap performance while considering the
// entry's enabled status.
func (r Race) FastestLap() []Entry {
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

		if sorted[i].Best < sorted[j].Best {
			return true
		} else {
			return false
		}
	})
	return sorted
}
