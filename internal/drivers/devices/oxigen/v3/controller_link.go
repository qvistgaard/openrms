package v3

import (
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/rs/zerolog"
	"time"
)

type controllerLink struct {
	id     types.CarId
	expire chan<- types.CarId
	renew  chan bool
	logger zerolog.Logger
}

func (l *controllerLink) timeout() {
	idleDuration := 15 * time.Second
	idleDelay := time.NewTimer(idleDuration)
	l.logger.Info().Int("id", int(l.id)).Msg("new controller link detected")
	for {
		idleDelay.Reset(idleDuration)

		select {
		case <-idleDelay.C:
			l.expire <- l.id
			l.logger.Warn().Int("id", int(l.id)).Msg("Controller link timed out")
			return
		case <-l.renew:

		}
	}
}
