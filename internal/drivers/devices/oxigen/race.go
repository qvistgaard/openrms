package oxigen

import log "github.com/sirupsen/logrus"

const (
	RaceUnknownByte           = 0x00
	RaceStoppedByte           = 0x01
	RaceRunningByte           = 0x03
	RacePausedByte            = 0x04
	RaceFlaggedLcEnabledByte  = 0x05
	RaceFlaggedLcDisabledByte = 0x15
)

type Race struct {
	status byte
}

func NewRace() *Race {
	return &Race{
		status: RaceUnknownByte,
	}
}

func (r *Race) Start() {
	log.WithField("drivers", "oxigen").
		WithField("race-status", "start").
		Info("Race status changed")
	r.status = RaceRunningByte
}

func (r *Race) Flag() {
	r.status = RaceFlaggedLcEnabledByte
}

func (r *Race) Pause() {
	log.WithField("drivers", "oxigen").
		WithField("race-status", "pause").
		Info("Race status changed")
	r.status = RacePausedByte
}

func (r *Race) Stop() {
	log.WithField("drivers", "oxigen").
		WithField("race-status", "stopepd").
		Info("Race status changed")
	r.status = RaceStoppedByte
}