package influxdb

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	api2 "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/qvistgaard/openrms/internal/telemetry"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"time"
)

type InfluxDB struct {
	client influxdb2.Client
	api    api2.WriteAPI
	/*	race   chan state.CourseState
		car    chan state.CarState*/
}

func (i *InfluxDB) Configure(observable rxgo.Observable) {
	observable.DistinctUntilChanged(func(ctx context.Context, i interface{}) (interface{}, error) {
		return i.(reactive.ValueChange).Value, nil
	}).DoOnNext(func(value interface{}) {
		i.processValueChange(value.(reactive.ValueChange))
	})
}

func (i *InfluxDB) processValueChange(change reactive.ValueChange) {
	if id, ok := change.Annotations[annotations.CarId]; ok {
		idInt := id.(types.Id)
		if field, ok := change.Annotations[annotations.CarValueFieldName]; ok {
			p := influxdb2.NewPointWithMeasurement("car")
			p.AddTag("id", idInt.String())
			p.SetTime(change.Timestamp)
			i.writePoint(p, field.(string), change.Value)
		}
	}
}

func (i *InfluxDB) writePoint(p *write.Point, s string, value interface{}) {
	i.processStateValue(p, s, value)
	i.api.WritePoint(p)
}

func (i InfluxDB) Process() {
	defer func() {
		log.Fatal("Influxdb process died")
	}()
	log.Info("started influxdb post processor.")
	for {
		select {
		case <-time.After(1 * time.Second):
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			i.api.WritePoint(
				influxdb2.NewPointWithMeasurement("memory").
					AddField("alloc", m.Alloc).
					AddField("total", m.TotalAlloc).
					AddField("sys", m.Alloc))

			i.api.WritePoint(
				influxdb2.NewPointWithMeasurement("gc").
					AddField("count", m.NumGC))
		}
	}
	log.Warn("influxdb processors stopped")
}

func (i *InfluxDB) processStateValue(p *write.Point, n string, v interface{}) {
	valueOf := reflect.ValueOf(v)
	kind := valueOf.Kind()
	if kind == reflect.Uint8 || kind == reflect.Uint {
		p.AddField(n, valueOf.Uint())
	} else if kind == reflect.Int64 {
		p.AddField(n, valueOf.Int())
	} else if kind == reflect.Uint16 {
		p.AddField(n, valueOf.Uint())
	} else if kind == reflect.String {
		p.AddField(n, valueOf.String())
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		p.AddField(n, valueOf.Float())
	} else if kind == reflect.Bool {
		p.AddField(n, valueOf.Bool())
	} else if kind == reflect.Ptr {
		if tv, ok := v.(telemetry.Metrics); ok {
			for _, va := range tv.Metrics() {
				i.processStateValue(p, va.Name, va.Value)
			}
		}
	} else if kind == reflect.Struct {
		if tv, ok := v.(*telemetry.Metrics); ok {
			for _, va := range (*tv).Metrics() {
				i.processStateValue(p, va.Name, va.Value)
			}
		} else if tv, ok := v.(telemetry.Metrics); ok {
			for _, va := range tv.Metrics() {
				i.processStateValue(p, va.Name, va.Value)
			}
		} else {
			log.WithField("state", n).
				WithField("value", v).
				WithField("kind", kind).
				WithField("type", fmt.Sprintf("%T", v)).
				Warn("Unsupported telemtry format")
		}
	} else {
		log.WithField("state", n).
			WithField("value", v).
			WithField("kind", kind).
			WithField("type", fmt.Sprintf("%T", v)).
			Warn("Unsupported telemtry format")
	}

}

/*
func (i *InfluxDB) CarChannel() chan<- state.CarState {
	return i.car
}

func (i *InfluxDB) RaceChannel() chan<- state.CourseState {
	return i.race
}
*/
