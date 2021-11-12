package car

import "github.com/qvistgaard/openrms/internal/state/rx/car"

type Rule interface {
	InitializeCarState(car *car.Car)
}
