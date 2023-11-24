package limbmode

import (
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
)

type Plugin struct {
	state     map[types.Id]observable.Observable[bool]
	carConfig map[types.Id]*LimbModeConfig
	config    *Config
}

func New(config *Config) (*Plugin, error) {
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

	return &Plugin{
		config:    config,
		carConfig: carConfig,
		state:     make(map[types.Id]observable.Observable[bool]),
	}, nil
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()

	if _, ok := p.carConfig[carId]; !ok {
		p.carConfig[carId] = p.config.Car.Defaults.LimbMode
	}

	p.state[carId] = observable.Create(false)
	p.state[carId].RegisterObserver(func(b bool, annotations observable.Annotations) {
		car.MaxSpeed().Update()
	})

	car.MaxSpeed().Modifier(func(u uint8) (uint8, bool) {
		return *p.carConfig[carId].MaxSpeed, p.state[carId].Get()
	}, 1)

	car.Pit().RegisterObserver(func(b bool, a observable.Annotations) {
		p.state[carId].Set(false)
	})
}

func (p *Plugin) InitializeCar(c *car.Car) {

}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.RaceStatus, annotations observable.Annotations) {

	})
}

func (p *Plugin) LimbMode(carId types.Id) observable.Observable[bool] {
	return p.state[carId]
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "limb-mode"
}
