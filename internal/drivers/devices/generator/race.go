package generator

import (
	"github.com/qvistgaard/openrms/internal/state/race"
	log "github.com/sirupsen/logrus"
	"time"
)

type Race struct {
	raceStatus race.Status
	laps       uint32
	raceStart  time.Time
}

func NewRace() *Race {
	return &Race{}
}

func (r *Race) Start() {
	log.Info("Race started")
	r.raceStatus = race.Running
	r.raceStart = time.Now()
	r.laps = 0
}

func (r *Race) Flag() {
	log.Info("Race Flagged")
	r.raceStatus = race.Flagged

}

func (r *Race) Pause() {
	log.Info("Race paused")
	r.raceStatus = race.Paused

}

func (r *Race) Stop() {
	log.Info("Race stopped")
	r.raceStatus = race.Stopped
}
