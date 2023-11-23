package track

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/observable"
)

type Track struct {
	implementer implement.Implementer
	maxSpeed    observable.Observable[uint8]
}

func New(c Config, driver implement.Implementer) (*Track, error) {
	var o implement.PitLaneLapCounting
	if c.Track.PitLane.LapCounting.OnEntry {
		o = implement.LapCountingOnEntry
	} else {
		o = implement.LapCountingOnExit
	}

	driver.Track().PitLane().LapCounting(c.Track.PitLane.LapCounting.Enabled, o)
	driver.Track().MaxSpeed(c.Track.MaxSpeed)

	return &Track{
		implementer: driver,
		maxSpeed:    observable.Create(uint8(0)),
	}, nil
}

func (t *Track) MaxSpeed() observable.Observable[uint8] {
	return t.maxSpeed
}

func (t *Track) Init(ctx context.Context) {
	t.maxSpeed.RegisterObserver(func(u uint8, annotations observable.Annotations) {
		t.implementer.Track().MaxSpeed(u)
	})
	// t.maxSpeed.Init(ctx)
}
