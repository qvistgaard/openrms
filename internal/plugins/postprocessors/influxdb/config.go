package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/qvistgaard/openrms/internal/state"
	"gopkg.in/yaml.v2"
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

func CreateFromConfig(config []byte) (*InfluxDB, error) {
	c := &Config{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}

	i := new(InfluxDB)
	db := c.Postprocessors.InfluxDB
	if db.BatchSize == 0 {
		db.BatchSize = 100
	}
	i.client = influxdb2.NewClientWithOptions(db.Url, db.AuthToken, influxdb2.DefaultOptions().SetBatchSize(db.BatchSize))
	i.api = i.client.WriteAPI(db.Organization, db.Bucket)
	i.race = make(chan state.RaceChanges, 1024)
	i.car = make(chan state.CarChanges, 1024)
	return i, nil
}
