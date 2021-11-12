package oxigen

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	log "github.com/sirupsen/logrus"
)

const (
	RaceRunningByte = 0x03
	RacePausedByte  = 0x03
	RaceStoppedByte = 0x01
)

type Race struct {
	status byte
}

func NewRace() *Race {
	return &Race{
		status: 0x00,
	}
}

func (r *Race) Status(status implement.RaceStatus) {
	switch status {
	case implement.RaceRunning:
		r.status = RaceRunningByte
	case implement.RacePaused:
		r.status = RacePausedByte
	case implement.RaceStopped:
		r.status = RaceStoppedByte
	default:
		log.Warn("race state not implemented")
	}
}

func (r *Race) Init(ctx context.Context, processor reactive.ValuePostProcessor) {
	panic("implement me")
}
