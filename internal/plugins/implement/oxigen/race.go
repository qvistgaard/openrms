package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
	log "github.com/sirupsen/logrus"
)

const (
	RaceUnknownByte = 0x00
	RaceStoppedByte = 0x01
	RaceRunningByte = 0x03
	RacePausedByte  = 0x04
)

type Race struct {
	status byte
}

func NewRace() *Race {
	return &Race{
		status: RaceUnknownByte,
	}
}

func (r *Race) Status(status implement.RaceStatus) {
	log.WithField("implement", "oxigen").
		WithField("state", status).
		Info("set race")
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
