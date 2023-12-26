package ontrack

import (
	"embed"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/commentary"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/utils"
	log "github.com/sirupsen/logrus"
)

//go:embed commentary/offtrack.txt
var announcements embed.FS

type Plugin struct {
	config           *Config
	flagPlugin       *flags.Plugin
	state            map[types.CarId]state
	flagId           int
	flag             flags.Flag
	commentaryPlugin *commentary.Plugin
}

type state struct {
	ontrack bool
	enabled bool
	inPit   bool
}

func New(c *Config, f *flags.Plugin, commentaryPlugin *commentary.Plugin) (*Plugin, error) {
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
		config:           c,
		flagPlugin:       f,
		commentaryPlugin: commentaryPlugin,
		flag:             flag,
		state:            make(map[types.CarId]state),
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

	car.Deslotted().RegisterObserver(func(b bool) {
		s := p.state[car.Id()]
		p.updateState(car.Id(), !b, s.inPit, s.enabled)

		if b && p.config.Plugin.OnTrack.Commentary {
			line, _ := utils.RandomLine(announcements, "commentary/offtrack.txt")
			template, err := utils.ProcessTemplate(line, car.TemplateData())
			if err == nil {
				p.commentaryPlugin.OptionalAnnouncement(template)
			}
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

func (p *Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p *Plugin) Priority() int {
	return 15
}

func (p *Plugin) Name() string {
	return "ontrack"
}
