package reactive

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type Factory struct {
	valuePostProcessor ValuePostProcessor
}

func NewFactory(valuePostProcessor ValuePostProcessor) *Factory {
	return &Factory{valuePostProcessor: valuePostProcessor}
}

func (f *Factory) NewDistinctValue(initial interface{}, annotations ...Annotations) Value {
	return newDistinctValueFunc(initial, defaultDistinctFunction(), f.valuePostProcessor, annotations...)
}

func (f *Factory) NewDistinctBoolean(initial bool, annotations ...Annotations) *Boolean {
	return &Boolean{f.NewDistinctValue(initial, annotations...)}
}

func (f *Factory) NewDuration(initial time.Duration, annotations ...Annotations) *Duration {
	return &Duration{f.NewDistinctValue(initial, annotations...)}
}

func (f *Factory) NewDistinctGauge(initial float64, annotations ...Annotations) *Gauge {
	return &Gauge{f.NewDistinctValue(initial, annotations...)}
}

func (f *Factory) NewLiter(initial types.Liter, annotations ...Annotations) *Liter {
	return &Liter{f.NewDistinctValue(initial, annotations...)}
}

func (f *Factory) NewDistinctPercent(initial types.Percent, annotations ...Annotations) *Percent {
	return &Percent{f.NewDistinctValue(initial, annotations...)}
}

func (f *Factory) NewPercent(initial types.Percent, annotations ...Annotations) *Percent {
	return &Percent{newValue(initial, f.valuePostProcessor, annotations...)}
}

func (f *Factory) NewDistinctLapNumber(annotations ...Annotations) *Lap {
	distinctValueFunc := func(ctx context.Context, i interface{}) (interface{}, error) {
		return i.(types.Lap).LapNumber, nil
	}
	return &Lap{newDistinctValueFunc(types.Lap{}, distinctValueFunc, f.valuePostProcessor, annotations...)}
}

func (f *Factory) NewDistictDrivers(initial types.Drivers, annotations ...Annotations) *Drivers {
	return &Drivers{f.NewDistinctValue(initial, annotations...)}
}
