package race

import (
	"github.com/gopxl/beep"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/plugins/race/sounds"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/system"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"time"
)

type Plugin struct {
	Duration     *time.Duration
	maxDuration  observable.Observable[time.Duration]
	Laps         *uint32
	status       race.Status
	confirmation *confirmation.Plugin
	confirmed    bool
	race         *race.Race
	started      bool
	fanfare      *streamer.Playback
	sound        *system.Sound
}

func New(r *race.Race, confirmationPlugin *confirmation.Plugin, sound *system.Sound) (*Plugin, error) {
	p := &Plugin{
		confirmation: confirmationPlugin,
		sound:        sound,
		maxDuration:  observable.Create(time.Duration(0)),
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
		if p.Duration != nil && p.status == race.Running {
			if *p.Duration <= duration {
				p.stop()
			}
		}
	})
	p.race.Laps().RegisterObserver(func(laps uint32) {
		if p.Laps != nil && *p.Laps <= laps && p.status == race.Running {
			p.stop()
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
			beeps := sounds.Beeps()
			p.sound.PlayEffect(beep.Seq(beeps, beep.Callback(func() {
				beeps.Close()
				p.confirmed = true
				p.race.Start()
				p.maxDuration.Set(*p.Duration)
				p.fanfare = nil
			})))
		}
	})
}

func (p *Plugin) stop() {
	p.race.Stop()
}

func (p *Plugin) MaxDuration() observable.Observable[time.Duration] {
	return p.maxDuration
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
