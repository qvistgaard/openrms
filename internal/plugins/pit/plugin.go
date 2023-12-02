package pit

import (
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	state    map[types.CarId]*state
	pitstops []SequencePlugin
}

func New(stops ...SequencePlugin) *Plugin {
	return &Plugin{
		state:    make(map[types.CarId]*state),
		pitstops: stops,
	}
}

type state struct {
	machine *stateless.StateMachine
	handler Handler
}

func (p *Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()
	handler := &DefaultHandler{
		car:      car,
		active:   observable.Create(false),
		current:  observable.Create(uint8(0)),
		maxSpeed: car.PitLaneMaxSpeed(),
	}
	p.state[carId] = &state{
		handler: handler,
		machine: machine(handler),
	}
	carState := p.state[carId]

	for _, ps := range p.pitstops {
		handler.sequences = append(handler.sequences, ps.ConfigurePitSequence(car.Id()))
	}

	car.PitLaneMaxSpeed().Modifier(func(u uint8) (uint8, bool) {
		return 0, handler.active.Get()
	}, 10000)
	car.Pit().RegisterObserver(func(b bool, annotations observable.Annotations) {
		var err error
		if !b {
			err = carState.machine.Fire(triggerCarExitedPitLane)
		} else {
			err = carState.machine.Fire(triggerCarEnteredPitLane)
		}
		if err != nil {
			log.Error(err)
		}
	})

	car.Controller().ButtonTrackCall().RegisterObserver(func(b bool, annotations observable.Annotations) {
		if b {
			err := carState.machine.Fire(triggerCarPitStopConfirmed)
			if err != nil {
				log.Error(err)
			}
		}
	})

	car.Controller().TriggerValue().RegisterObserver(func(u uint8, annotations observable.Annotations) {
		var err error
		if u == 0 {
			err = carState.machine.Fire(triggerCarStopped)
		} else {
			err = carState.machine.Fire(triggerCarMoving)
		}
		if err != nil {
			log.Error(err)
		}
	})

}

func (p *Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p *Plugin) Priority() int {
	return 100
}

func (p *Plugin) Name() string {
	return "pit"
}

func (p *Plugin) Active(car types.CarId) observable.Observable[bool] {
	return p.state[car].handler.Active()
}

func (p *Plugin) Current(car types.CarId) observable.Observable[uint8] {
	return p.state[car].handler.Current()
}
