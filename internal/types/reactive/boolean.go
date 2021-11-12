package reactive

type Boolean struct {
	Value
}

func NewBoolean(initial bool, annotations ...Annotations) *Boolean {
	return &Boolean{NewDistinctValue(initial, annotations...)}
}

func (p *Boolean) Set(state bool) {
	p.Value.Set(state)
}
