package generator

import "math/rand"

type controller struct {
}

func (c controller) BatteryWarning() bool {
	return false
}

func (c controller) Link() bool {
	return false
}

func (c controller) TrackCall() bool {
	return false
}

func (c controller) ArrowUp() bool {
	return false
}

func (c controller) ArrowDown() bool {
	return false
}

func (c controller) TriggerValue() float64 {
	return rand.Float64() * 100
}
