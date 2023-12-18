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
		number:      settings.Number,
	}

	// Initialize observable properties
	car.initObservableProperties(settings)

	// Register observers
	car.registerObservers()

	car.filters()
	return car
}

func (c *Car) initObservableProperties(settings *Settings) {
	c.maxBreaking = observable.Create(*settings.MaxBreaking).Filter(observable.DistinctComparableChange[uint8]())
	c.maxSpeed = observable.Create(*settings.MaxSpeed).Filter(observable.DistinctComparableChange[uint8]())
	c.minSpeed = observable.Create(*settings.MinSpeed).Filter(observable.DistinctComparableChange[uint8]())
	c.pitLaneMaxSpeed = observable.Create(*settings.PitLane.MaxSpeed).Filter(observable.DistinctComparableChange[uint8]())
	c.pit = observable.Create(false).Filter(observable.DistinctBooleanChange())
	c.deslotted = observable.Create(false).Filter(observable.DistinctBooleanChange())
	c.lastLap = observable.Create(types.Lap{})
	c.laps = observable.Create(uint32(0)).Filter(observable.DistinctComparableChange[uint32]())
	c.drivers = observable.Create(*settings.Drivers)
	c.team = observable.Create(*settings.Team).Filter(observable.DistinctComparableChange[string]())
	c.controller = controller.NewController()
	c.enabled = observable.Create(true).Filter(observable.DistinctBooleanChange())
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
	c.laps.Filter(observable.DistinctComparableChange[uint32]())
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
	enabled         observable.Observable[bool]
	number          *uint
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

func (c *Car) Number() uint {
	if c.number != nil {
		return *c.number
	}
	return uint(c.Id())
}

func (c *Car) Pit() observable.Observable[bool] {
	return c.pit
}

// Deslotted function
// Deprecated: Use OnTrack instead
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

func (c *Car) Enable() bool {
	return c.enabled.Set(true)
}

func (c *Car) Disable() bool {
	return c.enabled.Set(false)
}

func (c *Car) UpdateFromEvent(event drivers.Event) {
	switch e := event.(type) {
	case events.ControllerTriggerValueEvent:
		c.Controller().TriggerValue().Set(uint8(e.TriggerValue()))
	case events.Lap:
		if c.Laps().Set(e.Number()) {
			c.LastLap().Set(types.Lap{e.Number(), e.Time(), e.Recorded()})
		}
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
	c.enabled.Publish()
}

func (c *Car) Enabled() observable.Observable[bool] {
	return c.enabled
}
