package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"time"
)

type RaceStatus int

const (
	RaceStopped RaceStatus = iota
	RacePaused
	RaceRunning
	RaceFlagged
)

type RaceState struct {
	reactive.Value
}

func NewRaceState(initial RaceStatus, factory *reactive.Factory, annotations ...reactive.Annotations) *RaceState {
	return &RaceState{factory.NewDistinctValue(initial, annotations...)}
}

func (rs *RaceState) Set(value RaceStatus) {
	rs.Value.Set(value)
}

type Race struct {
	implementer implement.Implementer
	status      *RaceState
	raceTimer   *reactive.Duration

	raceStatus   RaceStatus
	raceDuration time.Duration
	raceStart    time.Time
}

func NewRace(implementer implement.Implementer, factory *reactive.Factory) *Race {
	return &Race{
		implementer: implementer,
		status: NewRaceState(RaceRunning, factory, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceStatus,
		}),
		raceTimer: factory.NewDuration(0, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceTimer,
		}),
		raceStart: time.Now(),
	}
}

func (r *Race) Start() {
	if r.raceStatus == RaceStopped {
		r.raceDuration = time.Second * 0
	}
	r.raceStart = time.Now()
	r.raceStatus = RaceRunning
	r.implementer.Race().Start()
	r.status.Set(r.raceStatus)
}

func (r *Race) Flag() {
	r.implementer.Race().Flag()
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.raceStatus = RaceFlagged
	r.status.Set(r.raceStatus)
}

func (r *Race) Stop() {
	r.implementer.Race().Stop()
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.raceStatus = RaceStopped
	r.status.Set(r.raceStatus)
}

func (r *Race) Pause() {
	r.implementer.Race().Pause()
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.raceStatus = RacePaused
	r.status.Set(r.raceStatus)
}

func (r *Race) Duration() time.Duration {
	if r.raceStatus == RaceRunning {
		return calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	} else {
		return r.raceDuration
	}
}

func calculateRaceDuration(duration time.Duration, startTime time.Time, currentTime time.Time) time.Duration {
	return duration + currentTime.Sub(startTime)
}

func (r *Race) CurrentState() RaceStatus {
	return r.raceStatus
}

func (r *Race) Status() *RaceState {
	return r.status
}

func (r *Race) Init(ctx context.Context) {
	r.status.Init(ctx)
	r.raceTimer.Init(ctx)
}
