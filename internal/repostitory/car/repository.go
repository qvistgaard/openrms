package car

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/rx/car"
	"github.com/qvistgaard/openrms/internal/types"
)

type Settings struct {
	Id       uint8
	MaxSpeed uint8 `yaml:"max-speed"`
}

type Repository interface {
	Get(id types.Id, ctx context.Context) (*car.Car, bool, bool)
	Exists(id types.Id) bool
	All() []*car.Car
}
