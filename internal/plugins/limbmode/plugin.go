package limbmode

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/plugins/commentary"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/utils"
)

//go:embed commentary/limbmode.txt
var announcements embed.FS

type Plugin struct {
	state      map[types.CarId]observable.Observable[bool]
	carConfig  map[types.CarId]*LimbModeConfig
	config     *Config
	commentary *commentary.Plugin
}

func New(config *Config, commentary *commentary.Plugin) (*Plugin, error) {
	carConfig := map[types.CarId]*LimbModeConfig{}
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
		config:     config,
		carConfig:  carConfig,
		commentary: commentary,
		state:      make(map[types.CarId]observable.Observable[bool]),
	}, nil
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()

	if _, ok := p.carConfig[carId]; !ok {
		p.carConfig[carId] = p.config.Car.Defaults.LimbMode
	}

	p.state[carId] = observable.Create(false).Filter(observable.DistinctBooleanChange())
	p.state[carId].RegisterObserver(func(b bool) {
		if b && p.config.Plugin.LimbMode.Commentary {
			line, err := utils.RandomLine(announcements, "commentary/limbmode.txt")
			if err == nil {
				template, _ := utils.ProcessTemplate(line, car.TemplateData())
				p.commentary.Announce(template)
			}
		}
		car.MaxSpeed().Update()
	})

	car.MaxSpeed().Modifier(func(u uint8) (uint8, bool) {
		return *p.carConfig[carId].MaxSpeed, p.state[carId].Get()
	}, 1)
}

func (p *Plugin) InitializeCar(_ *car.Car) {

}

func (p *Plugin) ConfigureRace(r *race.Race) {
	r.Status().RegisterObserver(func(status race.Status) {
		if status == race.Stopped {
			for _, o := range p.state {
				o.Set(false)
			}
		}
	})
}

func (p *Plugin) LimbMode(carId types.CarId) observable.Observable[bool] {
	return p.state[carId]
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "limb-mode"
}

func (p *Plugin) ConfigurePitSequence(carId types.CarId) pit.Sequence {
	return NewSequence(p.state[carId])
}

func (p *Plugin) Enabled() bool {
	return p.config.Plugin.LimbMode.Enabled
}
