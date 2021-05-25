package race

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Stages []StageConfig
}

func Create(config map[string]interface{}) *Rule {
	c := &Config{}
	err := mapstructure.Decode(config, c)
	if err != nil {
		log.Error(err)
	}

	for k, s := range c.Stages {
		if s.Type == nil {
			log.Fatalf("stage type is missing for stage: %d", k)
		}
		if s.Duration != nil {
			d, err := time.ParseDuration(*s.Duration)
			if err != nil {
				log.Fatalf("cant setup race duration. invalid duration: %s", *s.Duration)
			}
			c.Stages[k].duration = &d
		}
	}

	return &Rule{
		stages: c.Stages,
		ready:  make(map[state.CarId]*state.Car),
	}

}
