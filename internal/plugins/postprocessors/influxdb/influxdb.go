package influxdb

import (
	"github.com/influxdata/influxdb-client-go/v2"
	api2 "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type InfluxDB struct {
	client influxdb2.Client
	api    api2.WriteAPI
	race   chan state.CourseChanges
	car    chan state.CarChanges
}

func (i InfluxDB) Process() {
	log.Info("started influxdb post processor.")
	for {
		select {
		case car := <-i.car:
			p := influxdb2.NewPointWithMeasurement("car")
			for _, v := range car.Changes {
				p.AddField(v.Name, v.Value)
			}
			p.AddTag("id", strconv.Itoa(int(car.Car)))
			p.SetTime(car.Time)
			i.api.WritePoint(p)
		case race := <-i.race:
			p := influxdb2.NewPointWithMeasurement("race")
			for _, v := range race.Changes {
				p.AddField(v.Name, v.Value)
			}
			p.SetTime(race.Time)
			i.api.WritePoint(p)
		}
	}
	log.Warn("influxdb processors stopped")
}

func (i *InfluxDB) CarChannel() chan<- state.CarChanges {
	return i.car
}

func (i *InfluxDB) RaceChannel() chan<- state.CourseChanges {
	return i.race
}
