package reactive

import "github.com/qvistgaard/openrms/internal/types"

type Percent struct {
	Value
}

func NewPercent(initial types.Percent, annotations ...Annotations) *Percent {
	return &Percent{NewDistinctValue(initial, annotations...)}
}

func NewPercentAll(initial types.Percent, annotations ...Annotations) *Percent {
	return &Percent{NewValue(initial, annotations...)}
}

func (s Percent) Set(percent types.Percent) error {
	return s.Value.Set(percent)
}

func (s Percent) Get() types.Percent {
	return s.baseValue.(types.Percent)
}
