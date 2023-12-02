package race

import (
	"github.com/qvistgaard/openrms/internal/state/race"
	"time"
)

type Plugin struct {
	Duration *time.Duration
	Laps     *uint16
	status   race.Status
}

func New() (*Plugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Duration().RegisterObserver(func(duration time.Duration) {
		if p.Duration != nil && *p.Duration <= duration && p.status == race.Running {
			r.Stop()
		}
	})
	r.Laps().RegisterObserver(func(laps uint16) {
		if p.Laps != nil && *p.Laps <= laps && p.status == race.Running {
			r.Stop()
		}
	})
	r.Status().RegisterObserver(func(status race.Status) {
		p.status = status
	})
}

func (p *Plugin) Name() string {
	return "race"
}

func (p *Plugin) Priority() int {
	return 0
}
