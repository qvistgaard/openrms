package race

import (
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/state/race"
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

func New(r *race.Race, confirmationPlugin *confirmation.Plugin) (*Plugin, error) {
	p := &Plugin{
		confirmation: confirmationPlugin,
		race:         r,
	}

	p.initObservableProperties()
	p.registerObservers()

	return p, nil
}

func (p *Plugin) initObservableProperties() {

}

func (p *Plugin) registerObservers() {
	p.race.Duration().RegisterObserver(func(duration time.Duration) {
		if p.Duration != nil && *p.Duration <= duration && p.status == race.Running {
			p.race.Stop()
		}
	})
	p.race.Laps().RegisterObserver(func(laps uint32) {
		if p.Laps != nil && *p.Laps <= laps && p.status == race.Running {
			p.race.Stop()
		}
	})
	p.race.Status().Filter(func(status race.Status, status2 race.Status) bool {
		if status == race.Stopped && status2 == race.Running {
			// TODO fix this, make it wait
		}
		return false
	})
	p.race.Status().RegisterObserver(func(status race.Status) {
		p.status = status
	})
	p.confirmation.Confirmed().RegisterObserver(func(b bool) {
		if p.started && b {
			p.race.Start()
			p.started = false
		}
	})
}

func (p *Plugin) ConfigureRace(_ *race.Race) {
	// NOOP
}

/*
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

	func (p *Plugin) Race() *race.Race {
		return p.race
	}
*/
func (p *Plugin) Name() string {
	return "race"
}

func (p *Plugin) Priority() int {
	return 0
}
