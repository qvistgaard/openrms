package pit

import (
	"github.com/qvistgaard/openrms/internal/types"
)

type PitStop interface {
	Start() error
}

type Handler interface {
	PitStop
	Id() types.CarId
	OnCarStop(MachineTriggerFunc) error
	OnCarStart() error
	OnComplete() error
}
