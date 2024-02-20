package ontrack

import (
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/plugins/sound/system"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type Plugin struct {
	config     *Config
	flagPlugin *flags.Plugin
	state      map[types.CarId]state
	ontrack    map[types.CarId]observable.Observable[bool]
	flagId     int
	flag       flags.Flag
	sound      *system.Sound
	raceStatus race.Status
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.Status) {
		p.raceStatus = status
	})
}

type state struct {
	ontrack bool
	enabled bool
	inPit   bool
	cancel  chan bool
}

func New(c *Config, f *flags.Plugin, sound *system.Sound) (*Plugin, error) {
	var flag flags.Flag
	switch c.Plugin.OnTrack.Flag {
	case "green":
		flag = flags.Green
	case "yellow":
		flag = flags.Yellow
	case "red":
		flag = flags.Red
	default:
		return nil, errors.New("Invalid flagPlugin")
	}

	plugin := &Plugin{
		config:     c,
		flagPlugin: f,
		sound:      sound,
		flag:       flag,
		state:      make(map[types.CarId]state),
		ontrack:    make(map[types.CarId]observable.Observable[bool]),
	}

	if !f.Enabled() && c.Plugin.OnTrack.Enabled {
		log.WithField("plugin", plugin.Name()).
			Warn("flags plugin is not enabled, plugin will have no effect.")
	}

	return plugin, nil
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	p.state[car.Id()] = state{
		ontrack: true,
		inPit:   false,
		enabled: true,
	}

	p.ontrack[car.Id()] = observable.Create(true)

	car.Deslotted().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]

		if b {
			go func() {
				s.cancel = make(chan bool)
				select {
				case <-time.After(500 * time.Millisecond):
					p.updateState(car.Id(), !b, s.inPit, s.enabled)
				case <-s.cancel:
					close(s.cancel)
					s.cancel = nil
					return
				}
			}()
		} else {
			if s.cancel != nil {
				s.cancel <- true
			}
			p.updateState(car.Id(), !b, s.inPit, s.enabled)
		}
	})

	car.Pit().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), s.ontrack, b, s.enabled)
	})

	car.Enabled().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), s.ontrack, s.inPit, b)
	})
}

func (p *Plugin) updateState(id types.CarId, ontrack bool, inPit bool, enabled bool) {
	p.state[id] = state{
		ontrack: ontrack,
		inPit:   inPit,
		enabled: enabled,
	}
	p.ontrack[id].Set(ontrack)

	count := 0
	for _, s := range p.state {
		if !s.ontrack && s.enabled && !s.inPit {
			count = count + 1
		}
	}
	p.updateFlagPluginStatus(count)
}

func (p *Plugin) updateFlagPluginStatus(count int) {
	if count > 0 {
		log.WithField("deslotted", count).
			WithField("flag", p.flag).
			Info("Race flagged")
		if p.flagId < 0 {
			p.flagId = p.flagPlugin.Flag(p.flag)
		}
	} else {
		log.WithField("deslotted", count).
			Info("Flag cleared")
		p.flagPlugin.Clear(p.flagId)
		p.flagId = -1
	}
}

func (p *Plugin) Ontrack(id types.CarId) observable.Observable[bool] {
	return p.ontrack[id]
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
