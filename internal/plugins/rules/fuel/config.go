package fuel

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state/rx/rules"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Car struct {
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
}

type CarSettings struct {
	Id         *types.Id   `mapstructure:"id"`
	FuelConfig *FuelConfig `mapstructure:"fuel"`
}

type FuelConfig struct {
	TankSize     types.Liter          `mapstructure:"tank-size"`
	StartingFuel types.Liter          `mapstructure:"starting-fuel"`
	BurnRate     types.LiterPerSecond `mapstructure:"burn-rate"`
	FlowRate     types.LiterPerSecond `mapstructure:"flow-rate"`
}

func CreateFromConfig(applicationConfig *application.Config, rules rules.Rules) *Consumption {
	config := &Config{}
	err := mapstructure.Decode(applicationConfig, config)
	if err != nil {
		log.Error(err)
	}

	consumption := &Consumption{
		fuel:     make(map[types.Id]*reactive.Liter),
		state:    make(map[types.Id]*stateless.StateMachine),
		consumed: map[types.Id]*reactive.LiterSubtractModifier{},
		config:   config,
		rules:    rules,
	}
	return consumption
}
