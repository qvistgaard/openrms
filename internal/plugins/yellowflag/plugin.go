package yellowflag

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/race"
)

type Plugin struct {
	race *race.Race
}

func (p Plugin) ConfigureRace(r *race.Race) {
	//TODO implement me
	panic("implement me")
}

func New(race *race.Race) *Plugin {
	return &Plugin{race: race}
}

func (p Plugin) ConfigureCar(car *car.Car) {
	car.Deslotted().RegisterObserver(func(b bool) {
		if b {
			p.race.Pause()
		} else {
			p.race.Start()
		}
	})
}

func (p Plugin) InitializeCar(car *car.Car) {
	// NOOP
}

func (p Plugin) Priority() int {
	return 10
}

func (p Plugin) Name() string {
	return "yellow-flag"
}
