package generator

import (
	"github.com/qvistgaard/openrms/internal/state/race"
	log "github.com/sirupsen/logrus"
	"time"
)

type Race struct {
	raceStatus race.RaceStatus
	laps       uint16
	raceStart  time.Time
}

func NewRace() *Race {
	return &Race{}
}

func (r *Race) Start() {
	log.Info("Race started")
	r.raceStatus = race.RaceRunning
	r.raceStart = time.Now()
	r.laps = 0
}

func (r *Race) Flag() {
	log.Info("Race Flagged")
	r.raceStatus = race.RaceFlagged

}

func (r *Race) Pause() {
	log.Info("Race paused")
	r.raceStatus = race.RacePaused

}

func (r *Race) Stop() {
	log.Info("Race stopped")
	r.raceStatus = race.RaceStopped
}
