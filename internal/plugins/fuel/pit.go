package fuel

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Sequence struct {
	carState *state
}

func NewSequence(carState *state) *Sequence {
	return &Sequence{carState}
}

func (s *Sequence) Start() error {
	log.Info("Refuelling started.")
	full := false
	for !full {
		time.Sleep(250)
		s.carState.consumed, full = calculateRefuellingValue(s.carState.consumed, s.carState.config.FlowRate/4)
		s.carState.fuel.Update()
		log.WithField("fuel", s.carState.fuel.Get()).
			WithField("consumed", s.carState.consumed).
			WithField("full", full).
			Info("Refuelling.")

	}
	log.Info("Refuelling completed.")
	return nil
}

func calculateRefuellingValue(used float32, flowRate float32) (float32, bool) {
	liter := used - flowRate
	if liter <= 0 {
		return 0, true
	} else {
		return liter, false
	}
}
