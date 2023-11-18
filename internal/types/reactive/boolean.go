package reactive

type Boolean struct {
	Value
}

func (p *Boolean) Set(state bool) {
	p.Value.Set(state)
}
