package limbmode

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	MaxSpeed *state.Speed `mapstructure:"max-speed"`
}

func CreateFromConfig(config map[string]interface{}) *LimbMode {
	c := &Config{}
	err := mapstructure.Decode(config, c)
	if err != nil {
		log.Error(err)
	}
	if c.MaxSpeed == nil {
		speed := state.Speed(50)
		return &LimbMode{
			MaxSpeed: &speed,
		}
	}
	return &LimbMode{
		MaxSpeed: c.MaxSpeed,
	}

}
