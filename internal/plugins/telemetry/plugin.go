package telemetry

import (
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/ontrack"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type Plugin struct {
	listener       observable.Observable[Race]
	leader         observable.Observable[types.CarId]
	fastest        observable.Observable[types.CarId]
	telemetry      Race
	fuelPlugin     *fuel.Plugin
	limbModePlugin *limbmode.Plugin
	status         race.Status
	pitPlugin      *pit.Plugin
	ontrack        *ontrack.Plugin
}

func New(fuelPlugin *fuel.Plugin, limbModePlugin *limbmode.Plugin, pitPlugin *pit.Plugin, ontrack *ontrack.Plugin) *Plugin {
	return &Plugin{
		listener:       observable.Create(make(Race)),
		leader:         observable.Create(types.CarId(0)).Filter(observable.DistinctComparableChange[types.CarId]()),
		fastest:        observable.Create(types.CarId(0)).Filter(observable.DistinctComparableChange[types.CarId]()),
		telemetry:      make(Race),
		fuelPlugin:     fuelPlugin,
		limbModePlugin: limbModePlugin,
		pitPlugin:      pitPlugin,
		ontrack:        ontrack,
	}
}

func (p *Plugin) Priority() int {
	return 10000
}

func (p *Plugin) Name() string {
	return "telemetry"
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	id := car.Id()
	if entry, ok := p.telemetry[id]; !ok {
		entry = &Entry{
			car:    car,
			Id:     car.Id(),
			Number: car.Number(),
			Color:  car.Color(),
		}
		p.telemetry[id] = entry
		p.updateLeaderboard()
	}

	if f, err := p.fuelPlugin.Fuel(id); err == nil {
		f.RegisterObserver(func(f float32) {
			p.telemetry[id].Fuel = f
		})
	}
	car.Deslotted().RegisterObserver(func(b bool) {
		p.telemetry[id].Deslotted = b
	})

	car.Pit().RegisterObserver(func(b bool) {
		p.telemetry[id].InPit = b
	})

	car.MaxSpeed().RegisterObserver(func(u uint8) {
		p.telemetry[id].MaxSpeed = u
		p.updateLeaderboard()

	})
	car.MinSpeed().RegisterObserver(func(u uint8) {
		p.telemetry[id].MinSpeed = u
		p.updateLeaderboard()

	})
	car.PitLaneMaxSpeed().RegisterObserver(func(u uint8) {
		p.telemetry[id].MaxPitSpeed = u
	})
	if p.limbModePlugin.Enabled() {
		p.limbModePlugin.LimbMode(id).RegisterObserver(func(b bool) {
			p.telemetry[id].LimbMode = b
		})
	}

	car.Enabled().RegisterObserver(func(b bool) {
		p.telemetry[id].Enabled = b
		p.updateLeaderboard()

	})

	car.Team().RegisterObserver(func(s string) {
		p.telemetry[id].Team = s
		p.updateLeaderboard()

	})

	p.ontrack.Ontrack(id).RegisterObserver(func(b bool) {
		if b {
			p.telemetry[id].DeslotsLap = p.telemetry[id].DeslotsLap + 1
			p.telemetry[id].DeslotsTotal = p.telemetry[id].DeslotsTotal + 1
		}
	})

	car.LastLap().RegisterObserver(func(lap types.Lap) {
		p.telemetry[id].Laps = append(p.telemetry[id].Laps, lap)
		p.telemetry[id].Delta = time.Duration(lap.Time.Nanoseconds() - p.telemetry[id].Last.Time.Nanoseconds())
		p.telemetry[id].Last = lap
		p.telemetry[id].DeslotsLap = 0
		if p.telemetry[id].Best == 0 || p.telemetry[id].Last.Time < p.telemetry[id].Best {
			p.telemetry[id].Best = p.telemetry[id].Last.Time
		}
		p.updateLeaderboard()
	})

	if p.pitPlugin.Enabled() {
		p.pitPlugin.Active(id).RegisterObserver(func(b bool) {
			p.telemetry[id].PitStopActive = b
		})
		p.pitPlugin.Current(id).RegisterObserver(func(u uint8) {
			p.telemetry[id].PitStopSequence = u
		})
	}

}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.Status) {
		if status == race.Running && p.status == race.Stopped {
			for carId := range p.telemetry {
				p.telemetry[carId].Laps = nil
				p.telemetry[carId].Delta = 0
				p.telemetry[carId].Best = 0
				p.telemetry[carId].Last = types.Lap{}
			}
			p.status = race.Running
		}
		if status == race.Stopped && p.status == race.Running {
			p.status = status
			err := report(p.telemetry)
			if err != nil {
				log.Error(errors.WithMessage(err, "failed to write race report"))
			}
		}
	})
}

func (p *Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p *Plugin) RegisterObserver(observer observable.Observer[Race]) {
	p.listener.RegisterObserver(observer)
}

func (p *Plugin) Leader() observable.Observable[types.CarId] {
	return p.leader
}

func (p *Plugin) FastestLap() observable.Observable[types.CarId] {
	return p.fastest
}

func (p *Plugin) updateLeaderboard() {
	p.listener.Set(p.telemetry)

	go func() {
		leader := p.telemetry.Sort()
		if len(leader) > 0 {
			p.leader.Set(leader[0].Id)
		}
		fastest := p.telemetry.FastestLap()
		if len(fastest) > 0 {
			p.fastest.Set(fastest[0].Id)
		}
	}()
}
