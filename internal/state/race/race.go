package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"time"
)

type RaceStatus int

const (
	RaceStopped RaceStatus = iota
	RacePaused
	RaceRunning
	RaceFlagged
)

type Race struct {
	implementer implement.Implementer
	status      observable.Observable[RaceStatus]
	duration    observable.Observable[time.Duration]
	laps        observable.Observable[uint16]

	raceStatus   RaceStatus
	raceDuration time.Duration
	raceStart    time.Time
}

func New(_ Config, implementer implement.Implementer) (*Race, error) {
	return &Race{
		implementer: implementer,
		status:      observable.Create(RaceStopped, observable.Annotation{annotations.RaceValueFieldName, fields.RaceStatus}).Filter(filterRaceStatusChange()),
		duration:    observable.Create(time.Second*0, observable.Annotation{annotations.RaceValueFieldName, fields.RaceDuration}),                          // observable.Annotation{annotations.RaceValueFieldName, fields.RaceStatus},
		laps:        observable.Create(uint16(0), observable.Annotation{annotations.RaceValueFieldName, fields.Laps}).Filter(filterTotalLapsCountChange()), // observable.Annotation{annotations.RaceValueFieldName, fields.RaceStatus},
	}, nil
}

func (r *Race) Start() {
	r.raceStatus = RaceRunning
	r.implementer.Race().Start()
	r.status.Set(r.raceStatus)
}

func (r *Race) Flag() {
	r.implementer.Race().Flag()
	r.raceStatus = RaceFlagged
	r.status.Set(r.raceStatus)
}

func (r *Race) Stop() {
	r.implementer.Race().Stop()
	r.raceStatus = RaceStopped
	r.status.Set(r.raceStatus)
}

func (r *Race) Pause() {
	r.implementer.Race().Pause()
	r.raceStatus = RacePaused
	r.status.Set(r.raceStatus)
}

func (r *Race) Duration() observable.Observable[time.Duration] {
	return r.duration
}

func (r *Race) Laps() observable.Observable[uint16] {
	return r.laps
}

func (r *Race) Status() observable.Observable[RaceStatus] {
	return r.status
}

func (r *Race) Init(ctx context.Context) {

	// r.status.Init(ctx)
	// r.duration.Init(ctx)
}

func (r *Race) UpdateFromEvent(e implement.Event) {
	r.duration.Set(e.RaceTimer)
	r.laps.Set(e.Car.Lap.Number)
}
