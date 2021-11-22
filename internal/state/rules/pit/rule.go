package pit

import (
	"github.com/qvistgaard/openrms/internal/state/car"
)

type Rule interface {
	HandlePitStop(car *car.Car, cancel <-chan bool) bool
}
