package track

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/state/observable"
)

type Track struct {
	driver   drivers.Driver
	maxSpeed observable.Observable[uint8]
}

func New(c Config, di drivers.Driver) (*Track, error) {
	var o drivers.PitLaneLapCounting
	if c.Track.PitLane.LapCounting.OnEntry {
		o = drivers.LapCountingOnEntry
	} else {
		o = drivers.LapCountingOnExit
	}

	di.Track().PitLane().LapCounting(c.Track.PitLane.LapCounting.Enabled, o)
	di.Track().MaxSpeed(c.Track.MaxSpeed)

	t := &Track{
		driver:   di,
		maxSpeed: observable.Create(c.Track.MaxSpeed),
	}

	t.maxSpeed.RegisterObserver(func(u uint8, annotations observable.Annotations) {
		t.driver.Track().MaxSpeed(u)
	})
	return t, nil
}

func (t *Track) MaxSpeed() observable.Observable[uint8] {
	return t.maxSpeed
}

func (t *Track) Initialize() {
	t.maxSpeed.Publish()
}
