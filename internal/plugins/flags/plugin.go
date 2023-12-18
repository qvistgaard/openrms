package flags

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/track"
)

type Flag uint

const (
	Green Flag = iota
	Yellow
	Red
	White
	Checkered
)

type Plugin struct {
	race             *race.Race
	track            *track.Track
	config           *Config
	flagged          observable.Observable[Flag]
	activeFlags      map[int]Flag
	activeFlagConfig FlagConfig
	nextFlagId       int
}

func New(c *Config, t *track.Track, r *race.Race) (*Plugin, error) {
	plugin := &Plugin{
		track:      t,
		race:       r,
		config:     c,
		nextFlagId: 1,
		activeFlags: map[int]Flag{
			0: Green,
		},
	}
	plugin.initObservableProperties()
	plugin.registerObservers()

	return plugin, nil
}

func (p *Plugin) initObservableProperties() {
	p.flagged = observable.Create(Green).
		Filter(observable.DistinctComparableChange[Flag]()).
		Filter(p.flaggedIsPluginEnabledCondition)
	p.track.MaxSpeed().Modifier(p.trackMaxSpeedModifier, 10000)

}
func (p *Plugin) registerObservers() {
	p.flagged.RegisterObserver(p.handFlagUpdate)
}

func (p *Plugin) flaggedIsPluginEnabledCondition(_ Flag, _ Flag) bool {
	return p.Enabled() && p.race.Status().Get() != race.Stopped
}

func (p *Plugin) trackMaxSpeedModifier(_ uint8) (uint8, bool) {
	isActive := p.activeFlagConfig.MaxSpeed != nil
	if isActive {
		return *p.activeFlagConfig.MaxSpeed, isActive && p.Enabled()
	} else {
		return 0, false
	}
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "flags"
}

func (p *Plugin) Flagged() observable.Observable[Flag] {
	return p.flagged
}

func (p *Plugin) Flag(flag Flag) int {
	flagId := p.nextFlagId
	p.nextFlagId = flagId + 1

	p.activeFlags[flagId] = flag
	p.update()
	return flagId
}

func (p *Plugin) Clear(id int) {
	delete(p.activeFlags, id)
	p.update()
}

func (p *Plugin) HasActive(flag Flag) bool {
	for _, f := range p.activeFlags {
		if f == flag {
			return true
		}
	}
	return false
}

func (p *Plugin) update() {
	active := 0
	flag := Green
	for i, f := range p.activeFlags {
		if i > active && f > flag {
			active = i
		}
	}
	p.flagged.Set(p.activeFlags[active])
}

func (p *Plugin) handFlagUpdate(flag Flag) {
	p.activeFlagConfig = FlagConfig{}
	switch flag {
	case Green:
		p.race.Start()
	case Yellow:
		p.activeFlagConfig = p.config.Plugin.Flag.Yellow
	case Red:
		p.activeFlagConfig = p.config.Plugin.Flag.Red
	}

	if p.activeFlagConfig.Pause != nil && *p.activeFlagConfig.Pause {
		p.race.Pause()
	} else {
		p.track.MaxSpeed().Update()
	}
}

func (p *Plugin) Enabled() bool {
	return p.config.Plugin.Flag.Enabled
}
