package car

import (
	"github.com/divideandconquer/go-merge/merge"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/events"
	"github.com/qvistgaard/openrms/internal/state/controller"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func NewCar(implementer drivers.Driver, settings *Settings, defaults *Settings, id types.CarId) *Car {
	settings = merge.Merge(defaults, settings).(*Settings)

	car := &Car{
		implementer: implementer,
		id:          id,
	}

	// Initialize observable properties
	car.initObservableProperties(settings)

	// Register observers
	car.registerObservers()

	car.filters()
	return car
}

func (c *Car) initObservableProperties(settings *Settings) {
	c.maxBreaking = observable.Create(*settings.MaxBreaking)
	c.maxSpeed = observable.Create(*settings.MaxSpeed)
	c.minSpeed = observable.Create(*settings.MinSpeed)
	c.pitLaneMaxSpeed = observable.Create(*settings.PitLane.MaxSpeed)
	c.pit = observable.Create(false)
	c.deslotted = observable.Create(false)
	c.lastLap = observable.Create(types.Lap{})
	c.laps = observable.Create(uint32(0))
	c.drivers = observable.Create(*settings.Drivers)
	c.team = observable.Create(*settings.Team)
	c.controller = controller.NewController()
}

func (c *Car) registerObservers() {
	c.maxSpeed.RegisterObserver(func(u uint8) {
		c.implementer.Car(c.id).SetMaxSpeed(u)
	})
	c.minSpeed.RegisterObserver(func(u uint8) {
		c.implementer.Car(c.id).SetMinSpeed(u)
	})
	c.pitLaneMaxSpeed.RegisterObserver(func(u uint8) {
		c.implementer.Car(c.id).SetPitLaneMaxSpeed(u)
	})
	c.maxBreaking.RegisterObserver(func(u uint8) {
		c.implementer.Car(c.id).SetMaxBreaking(u)
	})
}

func (c *Car) filters() {
	c.maxSpeed.Filter(observable.DistinctPercentageChange())
	c.minSpeed.Filter(observable.DistinctPercentageChange())
	c.pitLaneMaxSpeed.Filter(observable.DistinctPercentageChange())
	c.maxBreaking.Filter(observable.DistinctPercentageChange())
	c.pit.Filter(observable.DistinctBooleanChange())
	c.deslotted.Filter(observable.DistinctBooleanChange())
	c.laps.Filter(observable.DistictComparableChange[uint32]())
}

type Car struct {
	id              types.CarId
	implementer     drivers.Driver
	controller      *controller.Controller
	pit             observable.Observable[bool]
	pitLaneMaxSpeed observable.Observable[uint8]
	maxSpeed        observable.Observable[uint8]
	minSpeed        observable.Observable[uint8]
	maxBreaking     observable.Observable[uint8]
	deslotted       observable.Observable[bool]
	laps            observable.Observable[uint32]
	lastLap         observable.Observable[types.Lap]
	drivers         observable.Observable[types.Drivers]
	team            observable.Observable[string]
}

func (c *Car) PitLaneMaxSpeed() observable.Observable[uint8] {
	return c.pitLaneMaxSpeed
}

func (c *Car) LastLap() observable.Observable[types.Lap] {
	return c.lastLap
}

func (c *Car) MaxSpeed() observable.Observable[uint8] {
	return c.maxSpeed
}

func (c *Car) MinSpeed() observable.Observable[uint8] {
	return c.minSpeed
}

func (c *Car) Controller() *controller.Controller {
	return c.controller
}

func (c *Car) Id() types.CarId {
	return c.id
}

func (c *Car) Pit() observable.Observable[bool] {
	return c.pit
}

func (c *Car) Deslotted() observable.Observable[bool] {
	return c.deslotted
}

func (c *Car) Laps() observable.Observable[uint32] {
	return c.laps
}

func (c *Car) Drivers() observable.Observable[types.Drivers] {
	return c.drivers
}

func (c *Car) Team() observable.Observable[string] {
	return c.team
}

func (c *Car) UpdateFromEvent(event drivers.Event) {
	switch e := event.(type) {
	case events.ControllerTriggerValueEvent:
		c.Controller().TriggerValue().Set(uint8(e.TriggerValue()))
	case events.Lap:
		c.Laps().Set(e.Number()) // get rid of this
		c.LastLap().Set(types.Lap{e.Number(), e.Time(), e.Recorded()})
		break
	case events.ControllerLinkEvent:
	case events.OnTrack:
	case events.ControllerTrackCallButton:
		c.Controller().ButtonTrackCall().Set(e.Pressed())
	case events.InPit:
		c.Pit().Set(e.InPit())
	case events.Deslotted:
		c.Deslotted().Set(e.Deslotted())
	default:
		log.WithField("package", reflect.TypeOf(e).Elem().PkgPath()).
			WithField("event", reflect.TypeOf(e).Elem().Name()).
			WithField("car", e.Car().Id()).Warn("Received unhandled event")
	}
}

func (c *Car) Initialize() {
	c.maxSpeed.Publish()
	c.pitLaneMaxSpeed.Publish()
	c.minSpeed.Publish()
	c.maxBreaking.Publish()
}
