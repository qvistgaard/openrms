package race

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/state/observable"
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
	implementer drivers.Driver
	status      observable.Observable[RaceStatus]
	duration    observable.Observable[time.Duration]
	laps        observable.Observable[uint16]

	raceStatus   RaceStatus
	raceDuration time.Duration
	raceStart    time.Time
}

func New(_ Config, implementer drivers.Driver) (*Race, error) {
	return &Race{
		implementer: implementer,
		status:      observable.Create(RaceStopped).Filter(filterRaceStatusChange()),
		duration:    observable.Create(time.Second * 0),
		laps:        observable.Create(uint16(0)).Filter(filterTotalLapsCountChange()), // observable.Annotation{annotations.RaceValueFieldName, fields.RaceStatus},
	}, nil
}

func (r *Race) Start() {
	r.raceStatus = RaceRunning
	r.raceStart = time.Now()
	if r.raceStatus == RaceStopped {
		r.raceDuration = time.Second * 0
	}
	r.implementer.Race().Start()
	r.status.Set(r.raceStatus)
}

func (r *Race) Flag() {
	r.implementer.Race().Flag()
	r.raceStatus = RaceFlagged
	r.status.Set(r.raceStatus)
}

func (r *Race) Stop() {
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.implementer.Race().Stop()
	r.raceStatus = RaceStopped
	r.status.Set(r.raceStatus)
}

func (r *Race) Pause() {
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
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

func (r *Race) UpdateFromEvent(e drivers.Event) {
	r.laps.Set(e.Car().Lap().Number())
	if r.raceStatus == RaceRunning {
		r.duration.Set(calculateRaceDuration(r.raceDuration, r.raceStart, time.Now()))
	}
}

func (r *Race) Initialize() {
	r.status.Publish()
}

func calculateRaceDuration(duration time.Duration, startTime time.Time, currentTime time.Time) time.Duration {
	return duration + currentTime.Sub(startTime)
}
