package reactive

type Gauge struct {
	Value
}

func NewGauge(initial float64, annotations ...Annotations) *Gauge {
	return &Gauge{NewDistinctValue(initial, annotations...)}
}

func (p *Gauge) Set(state float64) {
	p.Value.Set(state)
}
