package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	"net/http"
)

type RaceState struct {
	reactive.Value
}

func NewRaceState(initial implement.RaceStatus, annotations ...reactive.Annotations) *RaceState {
	return &RaceState{reactive.NewDistinctValue(initial, annotations...)}
}

func (p *RaceState) Set(value implement.RaceStatus) {
	p.Value.Set(value)
}

type Race struct {
	implementer implement.Implementer
	status      *RaceState
	raceTimer   *reactive.Duration
}

func (r *Race) RaceTimer() *reactive.Duration {
	return r.raceTimer
}

func NewRace(implementer implement.Implementer) *Race {
	return &Race{
		implementer: implementer,
		status: NewRaceState(implement.RaceStopped, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceStatus,
		}),
		raceTimer: reactive.NewDuration(0, reactive.Annotations{
			annotations.RaceValueFieldName: fields.RaceTimer,
		}),
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
	r.status.Init(ctx, postProcess)
	r.status.Update()

	r.RaceTimer().Init(ctx, postProcess)
}

func (r *Race) raceStatusChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		r.implementer.Race().Status(i.(implement.RaceStatus))
	})
}
