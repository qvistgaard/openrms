package leaderboard

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	"time"
)

type Leaderboard struct {
	listener  *reactive.RaceTelemetry
	telemetry types.RaceTelemetry
}

func New() *Leaderboard {
	return &Leaderboard{
		listener:  reactive.NewDistinctRaceTelemetry(),
		telemetry: make(types.RaceTelemetry),
	}
}

func (l *Leaderboard) Init(context context.Context) {
	l.listener.Init(context)
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
			var entry *types.RaceTelemetryEntry

			if entry, ok = l.telemetry[id]; !ok {
				entry = &types.RaceTelemetryEntry{
					Car: id,
				}
				l.telemetry[id] = entry
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
				// entry.PitState = change.Value.(types.CarPitState)
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
			case fields.MinSpeed:
				entry.MinSpeed = float64(change.Value.(types.Percent))
				l.updateLeaderboard()
			case fields.MaxTrackSpeed:
				entry.MaxSpeed = float64(change.Value.(types.Percent))
				l.updateLeaderboard()
			case fields.MaxPitSpeed:
				entry.MaxPitSpeed = float64(change.Value.(types.Percent))
				l.updateLeaderboard()
			case fields.Drivers:
				entry.Name = change.Value.(types.Drivers)[0].Name
				l.updateLeaderboard()
			}
		}
	}
	/*	if val, ok := change.Annotations[annotations.RaceValueFieldName]; ok {
		switch val {
		case fields.RaceTimer:
			l.raceTimer = change.Value.(time.Duration)
			l.updateLeaderboard()

		case fields.RaceStatus:
			l.raceStatus = change.Value.(implement.RaceStatus)
			l.updateLeaderboard()
		}
	}*/
}

func (l *Leaderboard) RegisterObserver(observer func(rxgo.Observable)) {
	l.listener.Value.RegisterObserver(observer)
}

func (l *Leaderboard) updateLeaderboard() {
	l.listener.Set(l.telemetry)
}

/*
func getDriverName(id types.Id, config *Config) string {
	for k, v := range config.Car.Cars {
		if types.Id(k) == id {
			return v.Drivers[0].Name
		}
	}
	return fmt.Sprintf("%s", getRandomDriver())
}

*/
