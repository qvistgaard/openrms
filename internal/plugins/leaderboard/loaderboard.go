package leaderboard

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type Plugin struct {
	listener  observable.Observable[types.RaceTelemetry]
	telemetry types.RaceTelemetry
}

func New() *Plugin {
	return &Plugin{
		listener:  observable.Create(make(types.RaceTelemetry)),
		telemetry: make(types.RaceTelemetry),
	}
}

func (p *Plugin) Priority() int {
	return 10000
}

func (p *Plugin) Name() string {
	return "leaderboard"
}

func (p *Plugin) ConfigureCar(car *car.Car) {

	id := car.Id()
	if entry, ok := p.telemetry[id]; !ok {
		entry = &types.RaceTelemetryEntry{
			Car: id,
		}
		p.telemetry[id] = entry
	}

	car.LastLap().RegisterObserver(func(lap types.Lap, a observable.Annotations) {
		p.telemetry[id].Laps = lap
		p.telemetry[id].Delta = time.Duration(lap.Time.Nanoseconds() - p.telemetry[id].Last.Nanoseconds())
		p.telemetry[id].Last = lap.Time
		if p.telemetry[id].Best == 0 || p.telemetry[id].Last < p.telemetry[id].Best {
			p.telemetry[id].Best = p.telemetry[id].Last
		}
		p.telemetry[id].Name = car.Team().Get()
		p.updateLeaderboard()
	})
}

func (p *Plugin) InitializeCarState(car *car.Car, ctx context.Context) {

}

/*
func (l *Plugin) processValueChange(change reactive.ValueChange) {
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
/*
}
*/

func (p *Plugin) RegisterObserver(observer observable.Observer[types.RaceTelemetry]) {
	p.listener.RegisterObserver(observer)
}

func (p *Plugin) updateLeaderboard() {
	p.listener.Set(p.telemetry)
}
