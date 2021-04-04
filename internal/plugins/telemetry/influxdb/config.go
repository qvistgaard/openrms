package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Telemetry struct {
		InfluxDB struct {
			Url          string `yaml:"url"`
			BatchSize    uint   `yaml:"batch-size"`
			AuthToken    string `yaml:"auth-token"`
			Organization string `yaml:"organization"`
			Bucket       string `yaml:"bucket"`
		} `yaml:"influxdb"`
	} `yaml:"telemetry"`
}

func CreateFromConfig(config []byte) (*InfluxDB, error) {
	c := &Config{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}

	i := new(InfluxDB)
	db := c.Telemetry.InfluxDB
	if db.BatchSize == 0 {
		db.BatchSize = 100
	}
	i.client = influxdb2.NewClientWithOptions(db.Url, db.AuthToken, influxdb2.DefaultOptions().SetBatchSize(db.BatchSize))
	i.api = i.client.WriteAPI(db.Organization, db.Bucket)
	return i, perr
}
