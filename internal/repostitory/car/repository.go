package car

import "github.com/qvistgaard/openrms/internal/state"

type Settings struct {
	Id       uint8
	MaxSpeed uint8 `yaml:"max-speed"`
}

type Repository interface {
	GetCarById(uint8 state.CarId) map[string]interface{}
	All() []*state.Car
}
