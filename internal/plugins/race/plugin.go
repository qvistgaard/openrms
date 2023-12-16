package race

import (
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/state/race"
	log "github.com/sirupsen/logrus"
	"time"
)

type Plugin struct {
	Duration     *time.Duration
	Laps         *uint32
	status       race.Status
	confirmation *confirmation.Plugin
	race         *race.Race
	started      bool
}

func (p *Plugin) Flag() {
	p.race.Flag()
}

func (p *Plugin) Pause() {
	p.race.Pause()
}

func (p *Plugin) Stop() {
	p.race.Stop()
}

func New(r *race.Race, confirmationPlugin *confirmation.Plugin) (*Plugin, error) {
	p := &Plugin{confirmation: confirmationPlugin, race: r}

	return p, nil
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Duration().RegisterObserver(func(duration time.Duration) {
		if p.Duration != nil && *p.Duration <= duration && p.status == race.Running {
			r.Stop()
		}
	})
	r.Laps().RegisterObserver(func(laps uint32) {
		if p.Laps != nil && *p.Laps <= laps && p.status == race.Running {
			r.Stop()
		}
	})
	r.Status().RegisterObserver(func(status race.Status) {
		p.status = status
	})

	p.confirmation.Confirmed().RegisterObserver(func(b bool) {
		if p.started && b {
			p.race.Start()
			p.started = false
		}
	})
}

func (p *Plugin) Start() {
	if !p.confirmation.Enabled() {
		p.race.Start()
		return
	}
	if !p.confirmation.Active().Get() {
		p.started = true
		err := p.confirmation.Activate()
		if err != nil {
			log.Error(err)
		}
	}
}

func (p *Plugin) Name() string {
	return "race"
}

func (p *Plugin) Priority() int {
	return 0
}
