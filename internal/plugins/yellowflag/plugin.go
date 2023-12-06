package yellowflag

import (
	race2 "github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/race"
)

type Plugin struct {
	race           *race.Race
	racePlugin     *race2.Plugin
	deslottedCount uint8
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	p.race = r
}

func New(r *race2.Plugin) *Plugin {
	return &Plugin{racePlugin: r, deslottedCount: 0}
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	car.Deslotted().RegisterObserver(func(b bool) {
		if b {
			p.deslottedCount = p.deslottedCount + 1
			p.race.Pause()
		} else {
			p.deslottedCount = p.deslottedCount - 1
			if p.deslottedCount <= 0 {
				p.deslottedCount = 0
				p.racePlugin.Start()
			}
		}
	})
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
