package reactive

import "time"

type Duration struct {
	Value
}

func NewDuration(initial time.Duration, annotations ...Annotations) *Duration {
	return &Duration{NewDistinctValue(initial, annotations...)}
}

func (p *Duration) Set(state time.Duration) {
	p.Value.Set(state)
}
