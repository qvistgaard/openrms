package reactive

import "github.com/qvistgaard/openrms/internal/types"

type Liter struct {
	Value
}

func NewLiter(initial types.Liter, annotations ...Annotations) *Liter {
	return &Liter{NewDistinctValue(initial, annotations...)}
}

func (s *Liter) Set(liter types.Liter) error {
	return s.Value.Set(liter)
}

func (s *Liter) Get() types.Liter {
	return s.value.(types.Liter)
}

type LiterModifierFunc func(types.Liter) types.Liter
type LiterModifier interface {
	Modify() LiterModifierFunc
}

func (s *Liter) Modifier(m LiterModifier) {
	s.Value.Modifier(func(i interface{}) interface{} {
		return m.Modify()(i.(types.Liter))
	})
}

type LiterSubtractModifier struct {
	Subtract types.Liter
}

func (l *LiterSubtractModifier) Modify() LiterModifierFunc {
	return func(liter types.Liter) types.Liter {
		t := liter - l.Subtract
		if t < 0 {
			return 0
		} else {
			return t
		}
	}
}
