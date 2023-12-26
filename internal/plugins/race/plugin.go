package race

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/plugins/commentary"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:embed commentary/start.txt
var announcements embed.FS

type Plugin struct {
	Duration     *time.Duration
	Laps         *uint32
	status       race.Status
	confirmation *confirmation.Plugin
	commentary   *commentary.Plugin
	confirmed    bool
	race         *race.Race
	started      bool
}

func New(r *race.Race, confirmationPlugin *confirmation.Plugin, commentary *commentary.Plugin) (*Plugin, error) {
	p := &Plugin{
		confirmation: confirmationPlugin,
		commentary:   commentary,
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
			if p.confirmed {
				p.confirmed = false
				return true
			}
			p.confirmation.Activate()
			return false
		}
		if status == race.Paused && status2 == race.Running {
			if p.confirmed {
				p.confirmed = false
				return true
			}
			p.confirmation.Activate()
			return false
		}
		return true
	})

	p.race.Status().RegisterObserver(func(status race.Status) {
		p.status = status
	})

	p.confirmation.Confirmed().RegisterObserver(func(b bool) {
		if b && (p.status == race.Stopped || p.status == race.Paused) {
			p.confirmed = true
			p.race.Start()
			line, err := utils.RandomLine(announcements, "commentary/start.txt")
			if err != nil {
				log.Error(err)
			} else {
				p.commentary.Announce(line)
			}
		}
	})
}

func (p *Plugin) ConfigureRace(_ *race.Race) {
	// NOOP
}

func (p *Plugin) Name() string {
	return "race"
}

func (p *Plugin) Priority() int {
	return 0
}
