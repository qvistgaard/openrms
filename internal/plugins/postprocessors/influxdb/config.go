package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

type Config struct {
	Postprocessors struct {
		InfluxDB struct {
			Url          string `yaml:"url"`
			BatchSize    uint   `yaml:"batch-size"`
			AuthToken    string `yaml:"auth-token"`
			Organization string `yaml:"organization"`
			Bucket       string `yaml:"bucket"`
		} `yaml:"influxdb"`
	} `yaml:"postprocessors"`
}

func CreateFromConfig(ctx *context.Context) (*InfluxDB, error) {
	c := &Config{}
	mapstructure.Decode(ctx.Config, c)
	i := new(InfluxDB)
	db := c.Postprocessors.InfluxDB
	if db.BatchSize == 0 {
		db.BatchSize = 100
	}
	i.client = influxdb2.NewClientWithOptions(db.Url, db.AuthToken, influxdb2.DefaultOptions().SetBatchSize(db.BatchSize))
	i.api = i.client.WriteAPI(db.Organization, db.Bucket)
	i.race = make(chan state.CourseChanges, 1024)
	i.car = make(chan state.CarChanges, 1024)
	return i, nil
}
