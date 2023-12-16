package ontrack

import (
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	flag   *flags.Plugin
	state  map[types.CarId]state
	flagId int
}

type state struct {
	ontrack bool
	enabled bool
}

func New(f *flags.Plugin) *Plugin {
	return &Plugin{
		flag:  f,
		state: make(map[types.CarId]state),
	}
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	p.state[car.Id()] = state{
		ontrack: true,
		enabled: true,
	}

	car.Deslotted().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), !b, s.enabled)
	})

	car.Enabled().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), s.ontrack, b)
	})
}

func (p *Plugin) updateState(id types.CarId, ontrack bool, enabled bool) {
	p.state[id] = state{
		ontrack: ontrack,
		enabled: enabled,
	}

	count := 0
	for _, s := range p.state {
		if !s.ontrack && s.enabled {
			count = count + 1
		}
	}
	if count > 0 {
		log.WithField("deslotted", count).
			Info("Yellow flagged")
		if p.flagId < 0 {
			p.flagId = p.flag.Flag(flags.Yellow)
		}
	} else {
		log.WithField("deslotted", count).
			Info("Flag cleared")
		p.flag.Clear(p.flagId)
		p.flagId = -1
	}
}

func (p *Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p *Plugin) Priority() int {
	return 15
}

func (p *Plugin) Name() string {
	return "ontrack"
}
