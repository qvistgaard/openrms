package car

import (
	"github.com/divideandconquer/go-merge/merge"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/state/controller"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

func NewCar(implementer drivers.Driver, settings *CarSettings, defaults *CarSettings, id types.Id) *Car {
	settings = merge.Merge(defaults, settings).(*CarSettings)

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

func (c *Car) initObservableProperties(settings *CarSettings) {
	c.maxBreaking = observable.Create(*settings.MaxBreaking)
	c.maxSpeed = observable.Create(*settings.MaxSpeed)
	c.minSpeed = observable.Create(*settings.MinSpeed)
	c.pitLaneMaxSpeed = observable.Create(*settings.PitLane.MaxSpeed)
	c.pit = observable.Create(false)
	c.deslotted = observable.Create(false)
	c.lastLapTime = observable.Create(0 * time.Second)
	c.lastLap = observable.Create(types.Lap{})
	c.laps = observable.Create(uint32(0))
	c.drivers = observable.Create(*settings.Drivers)
	c.team = observable.Create(*settings.Team)
	c.controller = controller.NewController()
}

func (c *Car) registerObservers() {
	c.maxSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).SetMaxSpeed(u)
	})
	c.minSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).SetMinSpeed(u)
	})
	c.pitLaneMaxSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).SetPitLaneMaxSpeed(u)
	})
	c.maxBreaking.RegisterObserver(func(u uint8, a observable.Annotations) {
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
	id              types.Id
	implementer     drivers.Driver
	controller      *controller.Controller
	pit             observable.Observable[bool]
	pitLaneMaxSpeed observable.Observable[uint8]
	maxSpeed        observable.Observable[uint8]
	minSpeed        observable.Observable[uint8]
	maxBreaking     observable.Observable[uint8]
	deslotted       observable.Observable[bool]
	lastLapTime     observable.Observable[time.Duration]
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

func (c *Car) Id() types.Id {
	return c.id
}

func (c *Car) Pit() observable.Observable[bool] {
	return c.pit
}

func (c *Car) Deslotted() observable.Observable[bool] {
	return c.deslotted
}

func (c *Car) LastLapTime() observable.Observable[time.Duration] {
	return c.lastLapTime
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

func (c *Car) UpdateFromEvent(e drivers.Event) {
	c.Pit().Set(e.Car().InPit())
	c.Deslotted().Set(e.Car().Deslotted())
	c.LastLapTime().Set(e.Car().Lap().Time())
	c.Laps().Set(uint32(e.Car().Lap().Number()))
	c.LastLap().Set(types.Lap{e.Car().Lap().Number(), e.Car().Lap().Time(), e.Car().Lap().Recorded()})
	c.Controller().ButtonTrackCall().Set(e.Car().Controller().TrackCall())
	c.Controller().TriggerValue().Set(uint8(e.Car().Controller().TriggerValue()))
}

func (c *Car) Initialize() {
	c.maxSpeed.Publish()
	c.pitLaneMaxSpeed.Publish()
	c.minSpeed.Publish()
	c.maxBreaking.Publish()
}
