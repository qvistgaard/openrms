package fuel

import (
	"time"
)

type Sequence struct {
	carState *state
}

func NewSequence(carState *state) *Sequence {
	return &Sequence{carState}
}

func (s *Sequence) Start() error {
	full := false
	for !full {
		s.carState.consumed, full = calculateRefuellingValue(s.carState.consumed, s.carState.config.FlowRate/4)
		s.carState.fuel.Update()
		time.Sleep(250)
	}
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
