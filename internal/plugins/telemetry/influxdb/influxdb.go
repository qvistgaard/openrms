package influxdb

import (
	"github.com/influxdata/influxdb-client-go/v2"
	api2 "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/qvistgaard/openrms/internal/telemetry"
)

func Connect() *InfluxDB {
	i := new(InfluxDB)
	i.client = influxdb2.NewClientWithOptions("http://localhost:8086", "_esnfjdFsjvkBkgc4S3ovqftb8RmmwAeGDBbt4fsY20EA9pFqeLaNYUkoH3sOyZ9jfTGcML-dB8AKgyj577Qhw==", influxdb2.DefaultOptions().SetBatchSize(100))
	i.api = i.client.WriteAPI("openrms", "openrms")
	return i
}

type InfluxDB struct {
	client influxdb2.Client
	api    api2.WriteAPI
}

func (i *InfluxDB) ProcessCar(t telemetry.Car) {
	p := influxdb2.NewPointWithMeasurement("openrms")
	for k, v := range t.Values {
		p.AddField(k, v)
	}
	p.AddTag("car", string(t.Car))
	p.SetTime(t.Time)
	i.api.WritePoint(p)
}
