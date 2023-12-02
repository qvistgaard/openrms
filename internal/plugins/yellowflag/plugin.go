package yellowflag

import (
	"github.com/qvistgaard/openrms/internal/state/car"
)

type Plugin struct {
}

func (p Plugin) ConfigureCar(car *car.Car) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) InitializeCar(car *car.Car) {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) Priority() int {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) Name() string {
	//TODO implement me
	panic("implement me")
}
