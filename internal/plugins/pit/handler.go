package pit

import (
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type DefaultHandler struct {
	cancel    chan bool
	car       *car.Car
	sequences []Sequence
	current   observable.Observable[uint8]
	active    observable.Observable[bool]
	maxSpeed  observable.Observable[uint8]
}

func (p *DefaultHandler) Active() observable.Observable[bool] {
	return p.active
}

func (p *DefaultHandler) Current() observable.Observable[uint8] {
	return p.current
}

func (p *DefaultHandler) OnComplete() error {
	log.WithField("car", p.Id()).Info("pit stop complete")
	return nil
}

func (p *DefaultHandler) OnCarStart() error {
	if p.cancel != nil && len(p.cancel) == 0 {
		p.cancel <- true
	}
	return nil
}

func (p *DefaultHandler) Id() types.CarId {
	return types.CarId(1)
}

func (p *DefaultHandler) OnCarStop(trigger MachineTriggerFunc) error {
	log.WithField("car", p.Id()).Info("car stopped inside pit lane")

	p.cancel = make(chan bool, 1)
	go func() {
		log.WithField("car", p.Id()).Info("waiting for automatic pit stop confirmation")
		select {
		case <-time.After(5 * time.Second):
			log.WithField("car", p.Id()).Info("pit stop automatically confirmed")
			err := trigger(triggerCarPitStopAutoConfirmed)
			if err != nil {
				log.Error(err)
			}
		case <-p.cancel:
			log.WithField("car", p.Id()).Info("pit stop wait cancelled")
		}
		defer func() { p.cancel = nil }()
	}()

	return nil
}

func (p *DefaultHandler) Start(trigger MachineTriggerFunc) error {
	go func() {
		p.active.Set(true)
		p.maxSpeed.Update()
		for i, sequence := range p.sequences {
			p.current.Set(uint8(i + 1))
			err := sequence.Start()
			if err != nil {
				log.Error(err, "pit stop sequence failed")
			}
		}
		p.active.Set(false)
		err := trigger(triggerCarPitStopComplete)
		if err != nil {
			log.Error(errors.WithMessage(err, "pit stop completion failed"))
			return
		}
	}()
	return nil
}
