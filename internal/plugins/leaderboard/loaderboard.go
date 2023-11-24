package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type Plugin struct {
	listener       observable.Observable[types.RaceTelemetry]
	telemetry      types.RaceTelemetry
	fuelPlugin     *fuel.Plugin
	limbModePlugin *limbmode.Plugin
}

func New(fuelPlugin *fuel.Plugin, limbModePlugin *limbmode.Plugin) *Plugin {
	return &Plugin{
		listener:       observable.Create(make(types.RaceTelemetry)),
		telemetry:      make(types.RaceTelemetry),
		fuelPlugin:     fuelPlugin,
		limbModePlugin: limbModePlugin,
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

	p.fuelPlugin.Fuel(id).RegisterObserver(func(f float32, annotations observable.Annotations) {
		p.telemetry[id].Fuel = f
	})

	car.Deslotted().RegisterObserver(func(b bool, annotations observable.Annotations) {
		p.telemetry[id].Deslotted = b
	})

	car.Pit().RegisterObserver(func(b bool, annotations observable.Annotations) {
		p.telemetry[id].InPit = b
	})

	car.MaxSpeed().RegisterObserver(func(u uint8, annotations observable.Annotations) {
		p.telemetry[id].MaxSpeed = u
	})
	p.limbModePlugin.LimbMode(id).RegisterObserver(func(b bool, annotations observable.Annotations) {
		p.telemetry[id].LimbMode = b

	})

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

func (p *Plugin) InitializeCar(c *car.Car) {

}

func (p *Plugin) RegisterObserver(observer observable.Observer[types.RaceTelemetry]) {
	p.listener.RegisterObserver(observer)
}

func (p *Plugin) updateLeaderboard() {
	p.listener.Set(p.telemetry)
}
