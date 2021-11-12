package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/qvistgaard/openrms/internal/webserver"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

const (
	CarLastLaps     = "car-last-laps"
	CarPosition     = "car-leaderboard-position"
	RaceLeaderboard = "race-leaderboard"
)

type Leaderboard struct {
	entries        map[types.Id]*BoardEntry
	webserver      webserver.WebServer
	lapUpdated     map[types.Id]*Lap
	emitterChannel chan rxgo.Item
	emmitter       rxgo.Disposed
}

func NewLeaderboard(ctx *application.Context) *Leaderboard {
	l := &Leaderboard{
		webserver:      ctx.Webserver,
		entries:        map[types.Id]*BoardEntry{},
		lapUpdated:     map[types.Id]*Lap{},
		emitterChannel: make(chan rxgo.Item),
	}
	l.emmitter = rxgo.FromChannel(l.emitterChannel).
		WindowWithTime(rxgo.WithDuration(1000 * time.Millisecond)).
		DoOnNext(func(i interface{}) {
			log.Info("recieved")
			o := i.(rxgo.Observable)
			last, err := o.Last().Get()
			if err == nil {
				l.webserver.PublishEvent(webserver.Event{
					Name:    "leaderboard-updated",
					Content: last.V.([]BoardEntry),
				})
			}
		})
	return l
}

type Lap struct {
	time   *time.Duration
	number *float64
}

type BoardEntry struct {
	Car       types.Id      `json:"car"`
	Laps      float64       `json:"lap"`
	Delta     time.Duration `json:"delta"`
	Best      time.Duration `json:"best"`
	Last      time.Duration `json:"last"`
	Deslotted bool          `json:"deslotted"`
	InPit     bool          `json:"in-pit"`
	Fuel      float64       `json:"fuel"`
	Name      string        `json:"name"`
}

func (l *Leaderboard) Configure(observable rxgo.Observable) {
	observable.DoOnNext(func(change interface{}) {
		l.processValueChange(change.(reactive.ValueChange))
	})
}

func (l *Leaderboard) processValueChange(change reactive.ValueChange) {
	if val, ok := change.Annotations[annotations.CarId]; ok {
		if field, ok := change.Annotations[annotations.CarValueFieldName]; ok {
			id := val.(types.Id)
			var entry *BoardEntry
			if entry, ok = l.entries[id]; !ok {
				entry = &BoardEntry{
					Car: id,
				}
				l.entries[id] = entry
				l.lapUpdated[id] = &Lap{}
			}

			switch field {
			case fields.LapTime:
				lapTime := change.Value.(time.Duration)
				entry.Delta = time.Duration(lapTime.Nanoseconds() - entry.Last.Nanoseconds())
				entry.Last = lapTime
				if entry.Best == 0 || entry.Last < entry.Best {
					entry.Best = entry.Last
				}
				l.lapUpdated[id].time = &lapTime
				l.updateLeaderboardIfLapChanged(id)
			case fields.Laps:
				entry.Laps = change.Value.(float64)
				l.lapUpdated[id].number = &entry.Laps
				l.updateLeaderboardIfLapChanged(id)

			case fields.InPit:
				entry.InPit = change.Value.(bool)
				l.updateLeaderboard()
			case fields.Fuel:
				t := change.Value.(types.Liter)
				entry.Fuel = t.ToFloat64()
				l.updateLeaderboard()
			case fields.Deslotted:
				entry.Deslotted = change.Value.(bool)
				l.updateLeaderboard()
			}
		}
	}
}

func (l *Leaderboard) updateLeaderboardIfLapChanged(id types.Id) {
	if l.lapUpdated[id].number != nil && l.lapUpdated[id].time != nil {
		l.lapUpdated[id] = &Lap{}
		l.updateLeaderboard()
	}
}

func (l *Leaderboard) updateLeaderboard() {
	sorted := make([]BoardEntry, 0, len(l.entries))

	for _, entry := range l.entries {
		sorted = append(sorted, *entry)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Laps > sorted[j].Laps {
			return true
		} else if sorted[i].Laps == sorted[j].Laps {
			// TODO: readd racetimer comparison
			/* if sorted[i].RaceTimer < sorted[j].RaceTimer {
				return true
			} */
			return false
		} else {
			return false
		}
	})
	l.emitterChannel <- rxgo.Of(sorted)

}

/*
   type Position uint

   type Leaderboard interface {
   	updateCar(id *state.Car, lap state.Lap) (Leaderboard, Position)
   }

   // Default Leaderboard implementation
   type Default struct {
   	entries []BoardEntry `json:"entries"`
   }

   func (l *Default) Compare(v state.ComparableChange) bool {
   	if c, ok := v.(*Default); ok {
   		return l == c
   	}
   	return true
   }

   func (l *Default) updateCar(car *state.Car, lap state.Lap) (Leaderboard, Position) {
   	r := &Default{entries: l.entries}
   	found := false
   	for k, v := range r.entries {
   		if v.Car == car.Id() {
   			r.entries[k].Lap = lap
   			r.entries[k].Delta = time.Duration(time.Duration(lap.LapTime).Nanoseconds() - time.Duration(r.entries[k].Best.LapTime).Nanoseconds())
   			if r.entries[k].Best.LapTime == 0 || lap.LapTime < r.entries[k].Best.LapTime {
   				r.entries[k].Best = lap
   			}
   			if settings, ok := car.GetSettings("driver").(map[interface{}]interface{}); ok {
   				r.entries[k].Name = settings["updateLeaderboard"].(string)
   			}
   			found = true

   			break
   		}
   	}
   	if !found {
   		r.entries = append(r.entries, BoardEntry{
   			Car: car.Id(),
   			Lap: lap,
   		})
   	}

   	sort.Slice(r.entries, func(i, j int) bool {
   		if r.entries[i].Lap.LapNumber > r.entries[j].Lap.LapNumber {
   			return true
   		} else if r.entries[i].Lap.LapNumber == r.entries[j].Lap.LapNumber {
   			if r.entries[i].Lap.RaceTimer < r.entries[j].Lap.RaceTimer {
   				return true
   			}
   			return false
   		} else {
   			return false
   		}
   	})
   	for p, v := range r.entries {
   		if v.Car == car.Id() {
   			return r, Position(p) + 1
   		}
   	}
   	return r, Position(0)
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
   		entries: []BoardEntry{},
   	})
   	b.Course = c
   	c.Subscribe(state.RaceStatus, b)
   }

   // TODO: test race reset of leader board
   // TODO: Reset leaderboard when race starts
   func (b *Rule) Notify(v *state.Value) {
   	if c, ok := v.Owner().(*state.Car); ok {
   		if rs, ok := b.Course.Get(state.RaceStatus).(uint8); !ok || rs != state.RaceStatusStopped {
   			if l, ok := v.Get().(*state.Lap); ok && v.Name() == state.CarLap {
   				last := c.Get(CarLastLaps).(LastLaps)
   				c.Set(CarLastLaps, last.update(*l))

   				leaderboard, i := b.Course.Get(RaceLeaderboard).(Leaderboard).updateCar(c, *l)
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
   							entries: []BoardEntry{},
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
*/
