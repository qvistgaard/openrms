package reactive

import "time"

type Duration struct {
	Value
}

func (p *Duration) Set(state time.Duration) {
	p.Value.Set(state)
}
