package generator

import (
	"github.com/qvistgaard/openrms/internal/state/race"
	log "github.com/sirupsen/logrus"
	"time"
)

type Race struct {
	raceStatus   race.RaceStatus
	raceDuration time.Duration
	raceStart    time.Time
	laps         uint16
}

func NewRace() *Race {
	return &Race{}
}

func (r *Race) Start() {
	log.Info("Race started")
	if r.raceStatus == race.RaceStopped {
		r.raceDuration = time.Second * 0
	}
	r.raceStart = time.Now()
	r.raceStatus = race.RaceRunning
	r.laps = 0
}

func (r *Race) Flag() {
	log.Info("Race Flagged")
	r.raceStatus = race.RaceFlagged
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())

}

func (r *Race) Pause() {
	log.Info("Race paused")
	r.raceStatus = race.RacePaused
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())

}

func (r *Race) Stop() {
	log.Info("Race stopped")
	r.raceStatus = race.RaceStopped
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
}

func calculateRaceDuration(duration time.Duration, startTime time.Time, currentTime time.Time) time.Duration {
	return duration + currentTime.Sub(startTime)
}
