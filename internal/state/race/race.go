package race

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"time"
)

type Status int

const (
	Stopped Status = iota
	Paused
	Running
	Flagged
)

type Race struct {
	implementer drivers.Driver
	status      observable.Observable[Status]
	duration    observable.Observable[time.Duration]
	laps        observable.Observable[uint32]

	raceStatus   Status
	raceDuration time.Duration
	raceStart    time.Time
}

func New(_ Config, implementer drivers.Driver) (*Race, error) {
	return &Race{
		implementer: implementer,
		status:      observable.Create(Stopped).Filter(filterRaceStatusChange()),
		duration:    observable.Create(time.Second * 0),
		laps:        observable.Create(uint32(0)).Filter(filterTotalLapsCountChange()), // observable.Annotation{annotations.RaceValueFieldName, fields.RaceStatus},
	}, nil
}

func (r *Race) Start() {
	if r.raceStatus == Stopped {
		r.raceDuration = time.Second * 0
	}
	r.raceStart = time.Now()
	r.raceStatus = Running
	r.implementer.Race().Start()
	r.status.Set(r.raceStatus)
}

func (r *Race) Flag() {
	r.implementer.Race().Flag()
	r.raceStatus = Flagged
	r.status.Set(r.raceStatus)
}

func (r *Race) Stop() {
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.implementer.Race().Stop()
	r.raceStatus = Stopped
	r.status.Set(r.raceStatus)
}

func (r *Race) Pause() {
	r.raceDuration = calculateRaceDuration(r.raceDuration, r.raceStart, time.Now())
	r.implementer.Race().Pause()
	r.raceStatus = Paused
	r.status.Set(r.raceStatus)
}

func (r *Race) Duration() observable.Observable[time.Duration] {
	return r.duration
}

func (r *Race) Laps() observable.Observable[uint32] {
	return r.laps
}

func (r *Race) Status() observable.Observable[Status] {
	return r.status
}

func (r *Race) UpdateFromEvent(event drivers.Event) {
	switch e := event.(type) {
	case events.Lap:
		r.laps.Set(e.Number())
	}

	if r.raceStatus == Running {
		r.duration.Set(calculateRaceDuration(r.raceDuration, r.raceStart, time.Now()))
	}
}

func (r *Race) Initialize() {
	r.status.Publish()
}

func calculateRaceDuration(duration time.Duration, startTime time.Time, currentTime time.Time) time.Duration {
	return duration + currentTime.Sub(startTime)
}
