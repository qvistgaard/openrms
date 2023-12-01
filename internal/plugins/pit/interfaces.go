package pit

import (
	"github.com/qvistgaard/openrms/internal/types"
)

type Stop interface {
	Start() error
	ConfigurePitStop(car Car)
}

type Car interface {
	Id() types.CarId
}

type Handler interface {
	Id() types.CarId
	OnCarStop(MachineTriggerFunc) error
	OnCarStart() error
	OnComplete() error
	Start(MachineTriggerFunc) error
}
