package v3

import (
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type controllerLink struct {
	id     types.CarId
	expire chan<- types.CarId
	renew  chan bool
}

func (l *controllerLink) timeout() {
	log.WithField("id", l.id).Info("new controller link detected")
	for {
		select {
		case <-time.After(1 * time.Second):
			l.expire <- l.id
			log.WithField("id", l.id).Warn("Controller link timed out")
			return
		case <-l.renew:

		}
	}
}
