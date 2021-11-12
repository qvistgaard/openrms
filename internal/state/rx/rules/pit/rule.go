package pit

import "github.com/qvistgaard/openrms/internal/state/rx/car"

type Rule interface {
	HandlePitStop(car *car.Car, cancel <-chan bool) bool
}
