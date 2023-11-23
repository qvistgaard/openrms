package reactive

/*
type Liter struct {
	Value
}

func (s *Liter) Set(liter types.Liter) error {
	return s.Value.Set(liter)
}

func (s *Liter) Get() types.Liter {
	return s.value.(types.Liter)
}

type LiterModifierFunc func(types.Liter) (types.Liter, bool)
type LiterModifier interface {
	Modify() LiterModifierFunc
}

func (s *Liter) Modifier(m LiterModifier, priority int) {
	s.Value.Modifier(func(i interface{}) (interface{}, bool) {
		return m.Modify()(i.(types.Liter))
	}, priority)
}

type LiterSubtractModifier struct {
	Subtract types.Liter
	Enabled  bool
}

func (l *LiterSubtractModifier) Modify() LiterModifierFunc {
	return func(liter types.Liter) (types.Liter, bool) {
		t := liter - l.Subtract
		if t < 0 {
			return 0, l.Enabled
		} else {
			return t, l.Enabled
		}
	}
}
*/
