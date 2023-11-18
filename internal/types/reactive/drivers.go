package reactive

import "github.com/qvistgaard/openrms/internal/types"

type Drivers struct {
	Value
}

func (p *Drivers) Set(state types.Drivers) {
	p.Value.Set(state)
}
