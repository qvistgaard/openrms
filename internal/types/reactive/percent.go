package reactive

import (
	"github.com/qvistgaard/openrms/internal/types"
)

type Percent struct {
	Value
}

func (s *Percent) Set(percent types.Percent) error {
	return s.Value.Set(percent)
}

func (s *Percent) Get() types.Percent {
	return s.baseValue.(types.Percent)
}

type PercentModifierFunc func(types.Percent) (types.Percent, bool)
type PercentModifier interface {
	Modify() PercentModifierFunc
}

func (s *Percent) Modifier(m PercentModifier, priority int) {
	s.Value.Modifier(func(i interface{}) (interface{}, bool) {
		return m.Modify()(i.(types.Percent))
	}, priority)
}

type PercentSubtractModifier struct {
	Subtract types.Percent
	Enabled  bool
}

func (p *PercentSubtractModifier) Modify() PercentModifierFunc {
	return func(percent types.Percent) (types.Percent, bool) {
		t := percent - p.Subtract
		if t < 0 {
			return 0, p.Enabled
		} else {
			return t, p.Enabled
		}
	}
}

type PercentAbsoluteModifier struct {
	Absolute  types.Percent
	Condition Condition
	Enabled   bool
}

func (p *PercentAbsoluteModifier) Modify() PercentModifierFunc {
	return func(percent types.Percent) (types.Percent, bool) {
		if p.Condition == IfGreaterThen {
			return p.Absolute, percent > p.Absolute && p.Enabled
		} else if p.Condition == IfLessThen {
			return p.Absolute, percent < p.Absolute && p.Enabled
		}
		return p.Absolute, p.Enabled
	}
}
