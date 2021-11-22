package track

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
)

type Track struct {
	implementer implement.Implementer
	maxSpeed    *reactive.Percent
}

func NewTrack(implementer implement.Implementer) *Track {
	return &Track{
		implementer: implementer,
		maxSpeed:    reactive.NewPercent(100),
	}
}

func (t *Track) MaxSpeed() *reactive.Percent {
	return t.maxSpeed
}

func (t *Track) Init(ctx context.Context, postProcess reactive.ValuePostProcessor) {
	t.maxSpeed.RegisterObserver(t.trackMaxSpeedChangeObserver)
	t.maxSpeed.Init(ctx, postProcess)
}

func (t *Track) trackMaxSpeedChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		t.implementer.Track().MaxSpeed(i.(types.Percent))
	})
}
