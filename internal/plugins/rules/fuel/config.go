package fuel

import (
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Fuel         *Liter
	StartingFuel *Liter          `mapstructure:"starting-fuel"`
	BurnRate     *LiterPerSecond `mapstructure:"burn-rate"`
	FlowRate     *LiterPerSecond `mapstructure:"flow-rate"`
}

func Create(config map[string]interface{}) *Consumption {
	c := &Config{}
	err := mapstructure.Decode(config, c)
	if err != nil {
		log.Error(err)
	}

	return &Consumption{config: c}

}
