package generator

import log "github.com/sirupsen/logrus"

type Race struct {
}

func NewRace() *Race {
	return &Race{}
}

func (r *Race) Start() {
	log.Info("Race started")
}

func (r *Race) Flag() {
	log.Info("Race Flagged")
}

func (r *Race) Pause() {
	log.Info("Race paused")
}

func (r *Race) Stop() {
	log.Info("Race stopped")
}
