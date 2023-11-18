package reactive

type Gauge struct {
	Value
}

func (p *Gauge) Set(state float64) {
	p.Value.Set(state)
}
