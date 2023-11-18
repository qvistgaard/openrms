package web

import (
	"fmt"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/qvistgaard/openrms/internal/webserver"
	"github.com/reactivex/rxgo/v2"
	"sort"
	"time"
)

type Leaderboard struct {
	entries        map[types.Id]*BoardEntry
	webserver      webserver.WebServer
	emitterChannel chan rxgo.Item
	emmitter       rxgo.Disposed
	raceTimer      time.Duration
	raceStatus     implement.RaceStatus
	config         *Config
}

type Event struct {
	RaceTimer   time.Duration
	RaceStatus  implement.RaceStatus
	Leaderboard []BoardEntry
}

func NewLeaderboard(ctx *application.Context, c *Config) *Leaderboard {
	l := &Leaderboard{
		webserver:      ctx.Webserver,
		entries:        map[types.Id]*BoardEntry{},
		emitterChannel: make(chan rxgo.Item),
		config:         c,
	}
	l.emmitter = rxgo.FromEventSource(l.emitterChannel).
		WindowWithTimeOrCount(rxgo.WithDuration(500*time.Millisecond), 1, rxgo.WithBufferedChannel(10)).
		DoOnNext(func(i interface{}) {
			o := i.(rxgo.Observable)
			last, err := o.Last().Get()
			if err == nil {
				sorted := last.V.([]BoardEntry)
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

				l.webserver.PublishEvent(webserver.Event{
					Name: "leaderboard-updated",
					Content: Event{
						RaceTimer:   l.raceTimer,
						RaceStatus:  l.raceStatus,
						Leaderboard: sorted,
					},
				})
			}
		})
	return l
}

/*type Lap struct {
	time    *time.Duration
	number  *float64
	lastLap types.Lap
}*/

type BoardEntry struct {
	Car       types.Id          `json:"car"`
	Laps      types.Lap         `json:"lap"`
	Delta     time.Duration     `json:"delta"`
	Best      time.Duration     `json:"best"`
	Last      time.Duration     `json:"last"`
	Deslotted bool              `json:"deslotted"`
	InPit     bool              `json:"in-pit"`
	Fuel      float64           `json:"fuel"`
	Name      string            `json:"name"`
	PitState  types.CarPitState `json:"pit-state"`
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
					Car:  id,
					Name: getDriverName(id, l.config),
				}
				l.entries[id] = entry
			}

			switch field {
			case fields.LastLap:
				lap := change.Value.(types.Lap)
				entry.Laps = lap
				l.updateLeaderboard()

			case fields.LapTime:
				lapTime := change.Value.(time.Duration)
				entry.Delta = time.Duration(lapTime.Nanoseconds() - entry.Last.Nanoseconds())
				entry.Last = lapTime
				if entry.Best == 0 || entry.Last < entry.Best {
					entry.Best = entry.Last
				}
				l.updateLeaderboard()
			case fields.PitState:
				entry.PitState = change.Value.(types.CarPitState)
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
	if val, ok := change.Annotations[annotations.RaceValueFieldName]; ok {
		switch val {
		case fields.RaceTimer:
			l.raceTimer = change.Value.(time.Duration)
			l.updateLeaderboard()

		case fields.RaceStatus:
			l.raceStatus = change.Value.(implement.RaceStatus)
			l.updateLeaderboard()
		}
	}
}

func getDriverName(id types.Id, config *Config) string {
	for k, v := range config.Car.Cars {
		if types.Id(k) == id {
			return v.Drivers[0].Name
		}
	}
	return fmt.Sprintf("%s", getRandomDriver())
}

func (l *Leaderboard) updateLeaderboard() {
	sorted := make([]BoardEntry, 0, len(l.entries))
	for _, entry := range l.entries {
		sorted = append(sorted, *entry)
	}
	l.emitterChannel <- rxgo.Of(sorted)
}
