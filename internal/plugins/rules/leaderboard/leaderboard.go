package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/telemetry"
	"sort"
)

const (
	CarLastLaps     = "car-last-laps"
	CarPosition     = "car-leaderboard-position"
	RaceLeaderboard = "race-leaderboard"
)

type Position uint

type Leaderboard interface {
	updateCar(id state.CarId, lap state.Lap) (Leaderboard, Position)
}

// Default Leaderboard implementation
type Default struct {
	Entries []BoardEntry `json:"entries"`
}

func (l *Default) Compare(v state.ComparableChange) bool {
	if c, ok := v.(*Default); ok {
		return l == c
	}
	return true
}

func (l *Default) updateCar(id state.CarId, lap state.Lap) (Leaderboard, Position) {
	r := &Default{Entries: l.Entries}
	found := false
	for k, v := range r.Entries {
		if v.Car == id {
			r.Entries[k].Lap = lap
			found = true
			break
		}
	}
	if !found {
		r.Entries = append(r.Entries, BoardEntry{
			Car: id,
			Lap: lap,
		})
	}

	sort.Slice(r.Entries, func(i, j int) bool {
		if r.Entries[i].Lap.LapNumber > r.Entries[j].Lap.LapNumber {
			return true
		} else if r.Entries[i].Lap.LapNumber == r.Entries[j].Lap.LapNumber {
			if r.Entries[i].Lap.RaceTimer < r.Entries[j].Lap.RaceTimer {
				return true
			}
			return false
		} else {
			return false
		}
	})
	for p, v := range r.Entries {
		if v.Car == id {
			return r, Position(p) + 1
		}
	}
	return r, Position(0)
}

type BoardEntry struct {
	Car state.CarId `json:"car"`
	Lap state.Lap   `json:"lap"`
}

type LastLaps interface {
	update(car state.Lap) *LastLapDefault
}
type LastLapDefault struct {
	Laps []state.Lap `json:"laps"`
}

type Rule struct {
	Course *state.Course
}

func (l *LastLapDefault) update(lap state.Lap) *LastLapDefault {
	return &LastLapDefault{
		Laps: append([]state.Lap{lap}, l.Laps[0:len(l.Laps)-1]...),
	}
}

func (l *LastLapDefault) Compare(v state.ComparableChange) bool {
	if c, ok := v.(*LastLapDefault); ok {
		return c.Laps[0].LapNumber == l.Laps[0].LapNumber &&
			c.Laps[0].LapTime == l.Laps[0].LapTime &&
			c.Laps[0].RaceTimer == l.Laps[0].RaceTimer
	}
	return true
}

func (b *Rule) InitializeCarState(car *state.Car) {
	car.Set(CarLastLaps, &LastLapDefault{
		Laps: make([]state.Lap, 5),
	})
	car.Subscribe(state.CarLap, b)
}

func (b *Rule) InitializeCourseState(c *state.Course) {
	c.Set(RaceLeaderboard, &Default{
		Entries: []BoardEntry{},
	})
	b.Course = c
	c.Subscribe(state.RaceStatus, b)
}

// TODO: test race reset of leader board
// TODO: Reset leaderboard when race starts
// TODO: add diff field
// TODO: Add best lap field
func (b *Rule) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		if rs, ok := b.Course.Get(state.RaceStatus).(uint8); !ok || rs != state.RaceStatusStopped {
			if l, ok := v.Get().(state.Lap); ok && v.Name() == state.CarLap {
				last := c.Get(CarLastLaps).(LastLaps)
				c.Set(CarLastLaps, last.update(l))

				leaderboard, i := b.Course.Get(RaceLeaderboard).(Leaderboard).updateCar(c.Id(), l)
				b.Course.Set(RaceLeaderboard, leaderboard)
				c.Set(CarPosition, i)
			}
		}
	} else if c, ok := v.Owner().(*state.Course); ok {
		if s, ok := v.Get().(uint8); ok && v.Name() == state.RaceStatus {
			if l, ok := v.GetPrevious().(uint8); ok {
				if s == state.RaceStatusRunning {
					if l != state.RaceStatusFlaggedLCDisabled &&
						l != state.RaceStatusFlaggedLCEnabled &&
						l != state.RaceStatusPaused {
						c.Set(RaceLeaderboard, &Default{
							Entries: []BoardEntry{},
						})
					}
				}
			}
		}
	}
}

func (l *Default) Metrics() []telemetry.Metric {
	return make([]telemetry.Metric, 0)
}

func (l *LastLapDefault) Metrics() []telemetry.Metric {
	return make([]telemetry.Metric, 0)
}
