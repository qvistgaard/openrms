package pit

import (
	"github.com/qmuntal/stateless"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	state map[types.CarId]*state
}

type state struct {
	machine *stateless.StateMachine
	handler Handler
}

func (p Plugin) ConfigureCar(car *car.Car) {
	carId := car.Id()
	handler := &DefaultHandler{car: car}
	p.state[carId] = &state{
		handler: handler,
		machine: machine(handler),
	}
	carState := p.state[carId]

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

func (p Plugin) InitializeCar(_ *car.Car) {
	// NOOP
}

func (p Plugin) Priority() int {
	return 100
}

func (p Plugin) Name() string {
	return "pit"
}
