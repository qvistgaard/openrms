package types

type Liter float64
type LiterPerSecond Liter

func (l Liter) ToFloat64() float64 {
	return float64(l)
}
