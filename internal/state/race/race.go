package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"time"
)

type RaceState struct {
	reactive.Value
}

func NewRaceState(initial implement.RaceStatus, factory *reactive.Factory, annotations ...reactive.Annotations) *RaceState {
	return &RaceState{factory.NewDistinctValue(initial, annotations...)}
}

func (p *RaceState) Set(value implement.RaceStatus) {
	p.Value.Set(value)
}

type Race struct {
	implementer       implement.Implementer
	status            *RaceState
	raceTimer         *reactive.Duration
	raceStart         time.Time
	raceStatusCurrent implement.RaceStatus
}

func (r *Race) UpdateTime() {
	if r.raceStatusCurrent == implement.RaceRunning {
		r.raceTimer.Set(time.Now().Sub(r.raceStart))
	}
}

func NewRace(implementer implement.Implementer, factory *reactive.Factory) *Race {
	return &Race{
		implementer: implementer,
		status: NewRaceState(implement.RaceRunning, factory, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceStatus,
		}),
		raceTimer: factory.NewDuration(0, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceTimer,
		}),
		raceStart: time.Now(),
	}
}

func (r *Race) Status() *RaceState {
	return r.status
}

func (r *Race) Init(ctx context.Context, postProcess reactive.ValuePostProcessor) {
	http.HandleFunc("/v1/race/start", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			r.status.Set(implement.RaceRunning)
			writer.WriteHeader(200)
			return
		}
		writer.WriteHeader(500)
	})
	http.HandleFunc("/v1/race/pause", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			r.status.Set(implement.RacePaused)
			writer.WriteHeader(200)
			return
		}
		writer.WriteHeader(500)
	})
	http.HandleFunc("/v1/race/stop", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			r.status.Set(implement.RaceStopped)
			writer.WriteHeader(200)
			return
		}
		writer.WriteHeader(500)
	})

	r.status.RegisterObserver(r.raceStatusChangeObserver)
	r.status.Init(ctx)
	r.status.Update()

	r.raceTimer.Init(ctx)
}

func (r *Race) raceStatusChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		status := i.(implement.RaceStatus)
		if status == implement.RaceRunning {
			r.raceStart = time.Now()
		}
		r.raceStatusCurrent = status
		r.implementer.Race().Status(status)
	})
}
