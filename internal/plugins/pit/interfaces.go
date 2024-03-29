package pit

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
)

type SequencePlugin interface {
	ConfigurePitSequence(types.CarId) Sequence
}

type Sequence interface {
	Start() error
}

type Car interface {
	Id() types.CarId
}

type Handler interface {
	Id() types.CarId
	OnCarStop(StartPitStop) error
	OnCarStart() error
	OnComplete() error
	Start(CompletePitStop, CancelPitStop) error
	Active() observable.Observable[bool]
	Current() observable.Observable[uint8]
}
