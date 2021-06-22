package influxdb

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	api2 "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/telemetry"
	log "github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strconv"
)

type InfluxDB struct {
	client influxdb2.Client
	api    api2.WriteAPI
	race   chan state.CourseState
	car    chan state.CarState
}

func (i InfluxDB) Process() {
	defer func() {
		log.Fatal("Influxdb process died")
	}()
	log.Info("started influxdb post processor.")
	for {
		select {
		case car := <-i.car:
			p := influxdb2.NewPointWithMeasurement("car")
			for _, v := range car.Changes {
				i.processStateValue(p, v.Name, v.Value)
			}
			p.AddTag("id", strconv.Itoa(int(car.Car)))
			p.SetTime(car.Time)
			i.api.WritePoint(p)
		case race := <-i.race:
			p := influxdb2.NewPointWithMeasurement("race")
			for _, v := range race.Changes {
				i.processStateValue(p, v.Name, v.Value)
			}
			p.SetTime(race.Time)
			i.api.WritePoint(p)
		}

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
	log.Warn("influxdb processors stopped")
}

func (i *InfluxDB) processStateValue(p *write.Point, n string, v interface{}) {
	valueOf := reflect.ValueOf(v)
	kind := valueOf.Kind()
	if kind == reflect.Uint8 || kind == reflect.Uint {
		p.AddField(n, valueOf.Uint())
	} else if kind == reflect.Int64 {
		p.AddField(n, valueOf.Int())
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

func (i *InfluxDB) CarChannel() chan<- state.CarState {
	return i.car
}

func (i *InfluxDB) RaceChannel() chan<- state.CourseState {
	return i.race
}
