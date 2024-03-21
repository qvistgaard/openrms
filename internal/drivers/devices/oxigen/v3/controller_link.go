package v3

import (
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/rs/zerolog/log"
	"time"
)

type controllerLink struct {
	id     types.CarId
	expire chan<- types.CarId
	renew  chan bool
}

func (l *controllerLink) timeout() {
	idleDuration := 15 * time.Second
	idleDelay := time.NewTimer(idleDuration)
	log.Info().Int("id", int(l.id)).Msg("new controller link detected")
	for {
		idleDelay.Reset(idleDuration)

		select {
		case <-idleDelay.C:
			l.expire <- l.id
			log.Warn().Int("id", int(l.id)).Msg("Controller link timed out")
			return
		case <-l.renew:

		}
	}
}
