package fuel

import (
	"github.com/qvistgaard/openrms/internal/types"
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
	TankSize     uint8                `mapstructure:"tank-size"`
	StartingFuel types.Liter          `mapstructure:"starting-fuel"`
	BurnRate     types.LiterPerSecond `mapstructure:"burn-rate"`
	FlowRate     types.LiterPerSecond `mapstructure:"flow-rate"`
}

/*func CreateFromConfig(applicationConfig *application.Config, rules rules.Rules) *Consumption {
	config := &Config{}
	err := mapstructure.Decode(applicationConfig, config)
	if err != nil {
		log.Error(err)
	}

	consumption := &Consumption{
		fuel:  make(map[types.Id]observable.Observable[float32]),
		state: make(map[types.Id]*stateless.StateMachine),
		// maxSpeed:   make(map[types.Id]*reactive.PercentSubtractModifier),
		fuelConfig: map[types.Id]*FuelConfig{},
		config:     config,
		rules:      rules,
	}
	return consumption
}
*/
