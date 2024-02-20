package confirmation

import (
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"slices"
	"time"
)

type Mode uint8

const (
	Timer Mode = 1 << iota
	TrackCall
)

type Plugin struct {
	mode          Mode
	timeout       time.Duration
	timerStart    time.Time
	confirmations map[types.CarId]bool
	confirmed     observable.Observable[bool]
	active        observable.Observable[bool]
	status        observable.Observable[Status]
	enabled       bool
	config        *Config
}

type Status struct {
	PendingConfirmations uint8
	TotalConfirmed       uint8
	RemainingTime        time.Duration
	TotalTime            time.Duration
}

func New(c *Config) (*Plugin, error) {
	p := &Plugin{
		active:    observable.Create(false).Filter(observable.DistinctBooleanChange()),
		confirmed: observable.Create(false).Filter(observable.DistinctBooleanChange()),
		status:    observable.Create(Status{}),
		timeout:   time.Second * 3,
		mode:      Timer,
		enabled:   c.Plugin.Confirmation.Enabled,
		config:    c,
	}

	if c.Plugin.Confirmation.Timeout != nil {
		p.timeout = *c.Plugin.Confirmation.Timeout
	}

	if c.Plugin.Confirmation.Modes != nil {
		if slices.Contains(c.Plugin.Confirmation.Modes, "timer") {
			p.mode = Timer ^ p.mode
		}
		if !slices.Contains(c.Plugin.Confirmation.Modes, "trackcall") {
			p.mode = TrackCall ^ p.mode
		}
	}

	p.active.RegisterObserver(func(b bool) {
		if !b {
			p.confirmed.Set(false)
		}
	})

	return p, nil
}

func (p *Plugin) Priority() int {
	return 0
}

func (p *Plugin) Name() string {
	return "confirmation"
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	if p.mode&TrackCall != 0 {
		car.Controller().ButtonTrackCall().RegisterObserver(func(b bool) {
			if b && p.active.Get() {
				status := p.status.Get()
				status.TotalConfirmed = status.TotalConfirmed + 1
				status.PendingConfirmations = status.PendingConfirmations - 1
				p.status.Set(status)
				p.confirmations[car.Id()] = true
				for _, b2 := range p.confirmations {
					if !b2 {
						return
					}
				}
				p.confirmed.Set(true)
				p.active.Set(false)
			}
		})
	}
}

func (p *Plugin) Activate() error {
	if p.active.Get() {
		return errors.New("Confirmation already in progress")
	}
	log.Info("Confirmation process started")

	p.confirmed.Set(false)
	p.active.Set(true)

	status := p.status.Get()

	if p.mode&Timer != 0 {
		p.timerStart = time.Now()

		go func() {
			for {
				status := p.status.Get()
				time.Sleep(100 * time.Millisecond)
				status.RemainingTime = p.timeout - time.Now().Sub(p.timerStart)
				status.TotalTime = p.timeout
				p.status.Set(status)
				if 0 >= status.RemainingTime {
					p.confirmed.Set(true)
					p.active.Set(false)

					log.WithField("mode", "timer").
						Info("Confirmation process completed")

					return
				}
			}
		}()
	}

	if p.mode&TrackCall != 0 {
		status.PendingConfirmations = uint8(len(p.confirmations))
		status.TotalConfirmed = 0
		for id := range p.confirmations {
			p.confirmations[id] = false
		}
	}
	p.status.Set(status)

	return nil
}

func (p *Plugin) Active() observable.Observable[bool] {
	return p.active
}

func (p *Plugin) Status() observable.Observable[Status] {
	return p.status
}

func (p *Plugin) Confirmed() observable.Observable[bool] {
	return p.confirmed
}

func (p *Plugin) Enabled() bool {
	return p.enabled
}

func (p *Plugin) InitializeCar(car *car.Car) {
	// NOOP
}
