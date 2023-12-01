package pit

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type DefaultHandler struct {
	cancel  chan bool
	car     *car.Car
	actions []*Stop
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
	log.Info("PIT RUNNING")

	return trigger(triggerCarPitStopComplete)
}
