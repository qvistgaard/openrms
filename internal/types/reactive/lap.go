package reactive

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
)

type Lap struct {
	Value
}

func NewLap(annotations ...Annotations) *Lap {
	distinctValueFunc := func(ctx context.Context, i interface{}) (interface{}, error) {
		return i.(types.Lap).LapNumber, nil
	}
	return &Lap{NewDistinctValueFunc(types.Lap{}, distinctValueFunc, annotations...)}
}

func (p *Lap) Set(state types.Lap) {
	p.Value.Set(state)
}
