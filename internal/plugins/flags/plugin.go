package flags

import (
	"github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/state/observable"
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
	race        *race.Plugin
	flagged     observable.Observable[Flag]
	activeFlags map[int]Flag
	nextFlagId  int
}

func New(r *race.Plugin) *Plugin {
	plugin := &Plugin{
		race:       r,
		nextFlagId: 1,
		activeFlags: map[int]Flag{
			0: Green,
		},
	}
	plugin.initObservableProperties()
	plugin.registerObservers()

	return plugin
}

func (p *Plugin) initObservableProperties() {
	p.flagged = observable.Create(Green).Filter(observable.DistinctComparableChange[Flag]())
}

func (p *Plugin) registerObservers() {
	p.flagged.RegisterObserver(p.handFlagUpdate)
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "yellow-flag"
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
	switch flag {
	case Green:
		p.race.Start()
	case Yellow:
		p.race.Pause()
	}
}
