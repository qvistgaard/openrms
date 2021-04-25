package car

import "github.com/qvistgaard/openrms/internal/state"

type Settings struct {
	Id       uint8
	MaxSpeed uint8 `yaml:"max-speed"`
}

type Repository interface {
	Get(id state.CarId) (*state.Car, bool)
	All() []*state.Car
}
