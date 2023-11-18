package reactive

import (
	"github.com/qvistgaard/openrms/internal/types"
)

type RaceTelemetry struct {
	Value
}

func (p *RaceTelemetry) Set(state types.RaceTelemetry) {
	p.Value.Set(state)
}

func NewDistinctRaceTelemetry(annotations ...Annotations) *RaceTelemetry {
	return &RaceTelemetry{
		newValue(make(types.RaceTelemetry), emptyValueProcessor(), annotations...),
	}
}
