package yellowflag

import (
	race2 "github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	race       *race.Race
	racePlugin *race2.Plugin
	state      map[types.CarId]state
	flagged    observable.Observable[bool]
}

type state struct {
	deslotted bool
	enabled   bool
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	p.race = r
}

func New(r *race2.Plugin) *Plugin {
	return &Plugin{
		racePlugin: r,
		state:      make(map[types.CarId]state),
		flagged:    observable.Create(false),
	}
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	p.state[car.Id()] = state{
		deslotted: false,
		enabled:   true,
	}

	car.Deslotted().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), b, s.enabled)
	})

	car.Enabled().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), s.deslotted, b)
	})
}

func (p *Plugin) updateState(id types.CarId, deslotted bool, enabled bool) {
	p.state[id] = state{
		deslotted: deslotted,
		enabled:   enabled,
	}

	deslottedCount := 0
	for _, s := range p.state {
		if s.deslotted && s.enabled {
			deslottedCount = deslottedCount + 1
		}
	}
	if deslottedCount > 0 {
		log.WithField("deslotted", deslottedCount).
			Info("Yellow flagged")
		p.race.Pause()
	} else {
		log.WithField("deslotted", deslottedCount).
			Info("Restarted race")
		p.racePlugin.Start()
	}
}

func (p *Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "yellow-flag"
}
