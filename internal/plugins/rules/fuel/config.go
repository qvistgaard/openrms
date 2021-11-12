package fuel

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Fuel         *types.Liter
	StartingFuel *types.Liter          `mapstructure:"starting-fuel"`
	BurnRate     *types.LiterPerSecond `mapstructure:"burn-rate"`
	FlowRate     *types.LiterPerSecond `mapstructure:"flow-rate"`
}

func Create(config map[string]interface{}, postprocessors *postprocess.PostProcess) *Consumption {
	c := &Config{}
	err := mapstructure.Decode(config, c)
	if err != nil {
		log.Error(err)
	}

	consumption := &Consumption{
		fuel:          make(map[types.Id]*reactive.Liter),
		state:         make(map[types.Id]*stateless.StateMachine),
		consumed:      map[types.Id]*reactive.LiterSubtractModifier{},
		config:        c,
		postprocessor: postprocessors,
	}
	return consumption
}
