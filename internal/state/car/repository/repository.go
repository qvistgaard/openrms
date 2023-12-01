package repository

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types"
)

type Settings struct {
	Id       uint8
	MaxSpeed uint8 `yaml:"max-speed"`
}

type Repository interface {
	Get(id types.CarId) (*car.Car, bool, bool)
}
