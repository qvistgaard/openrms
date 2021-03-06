package limbmode

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	log "github.com/sirupsen/logrus"
)

type LimbModeConfig struct {
	MaxSpeed *types.Percent `mapstructure:"max-speed"`
}

type CarSettings struct {
	Id       *types.Id       `mapstructure:"id"`
	LimbMode *LimbModeConfig `mapstructure:"limb-mode"`
}

type Config struct {
	Car struct {
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
}

func CreateFromConfig(applicationConfig *application.Config) *LimbMode {
	config := &Config{}

	err := mapstructure.Decode(applicationConfig, config)
	if err != nil {
		log.Error(err)
	}

	carConfig := map[types.Id]*LimbModeConfig{}
	for _, v := range config.Car.Cars {
		if v.LimbMode == nil {
			v.LimbMode = &LimbModeConfig{}
		}
		if v.LimbMode.MaxSpeed == nil {
			v.LimbMode.MaxSpeed = config.Car.Defaults.LimbMode.MaxSpeed
		}
		carConfig[*v.Id] = v.LimbMode
	}

	return &LimbMode{
		defaults:      config.Car.Defaults.LimbMode,
		config:        carConfig,
		state:         map[types.Id]*reactive.Boolean{},
		speedModifier: map[types.Id]*reactive.PercentAbsoluteModifier{},
	}
}
