package car

import (
	"context"
	"github.com/divideandconquer/go-merge/merge"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state/controller"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"time"
)

func NewCar(implementer implement.Implementer, settings *CarSettings, defaults *CarSettings, id types.Id) *Car {
	settings = merge.Merge(defaults, settings).(*CarSettings)

	annotations := []observable.Annotation{
		{annotations.CarId, id.String()},
	}

	car := &Car{
		implementer: implementer,
		id:          id,
	}

	// Initialize observable properties
	car.initObservableProperties(settings, annotations)

	// Register observers
	car.registerObservers()

	return car
}

func (c *Car) initObservableProperties(settings *CarSettings, a []observable.Annotation) {
	c.maxBreaking = observable.Create(*settings.MaxBreaking)
	c.maxSpeed = observable.Create(*settings.MaxSpeed, append(a, observable.Annotation{annotations.CarValueFieldName, fields.MaxTrackSpeed})...)
	c.minSpeed = observable.Create(*settings.MinSpeed, append(a, observable.Annotation{annotations.CarValueFieldName, fields.MinSpeed})...)
	c.pitLaneMaxSpeed = observable.Create(*settings.PitLane.MaxSpeed, append(a, observable.Annotation{annotations.CarValueFieldName, fields.MaxPitSpeed})...)
	c.pit = observable.Create(false, append(a, observable.Annotation{annotations.CarValueFieldName, fields.InPit})...).Filter(observable.DistinctBooleanChange())
	c.deslotted = observable.Create(false, append(a, observable.Annotation{annotations.CarValueFieldName, fields.Deslotted})...).Filter(observable.DistinctBooleanChange())
	c.lastLapTime = observable.Create(0*time.Second, observable.Annotation{annotations.CarValueFieldName, fields.LapTime})
	c.lastLap = observable.Create(types.Lap{}, append(a, observable.Annotation{annotations.CarValueFieldName, fields.LastLap})...)
	c.laps = observable.Create(uint32(0), append(a, observable.Annotation{annotations.CarValueFieldName, fields.Laps})...)
	c.drivers = observable.Create(*settings.Drivers, append(a, observable.Annotation{annotations.CarValueFieldName, fields.Drivers})...)
	c.team = observable.Create(*settings.Team, append(a, observable.Annotation{annotations.CarValueFieldName, fields.Drivers})...)
	c.controller = controller.NewController(a...)
}

func (c *Car) registerObservers() {
	c.maxSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).MaxSpeed(u)
	})
	c.minSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).MinSpeed(u)
	})
	c.pitLaneMaxSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).PitLaneMaxSpeed(u)
	})
	c.maxBreaking.RegisterObserver(func(u uint8, a observable.Annotations) {
		c.implementer.Car(c.id).PitLaneMaxSpeed(u)
	})
}

type Car struct {
	id              types.Id
	implementer     implement.Implementer
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

func (c *Car) UpdateFromEvent(e implement.Event) {
	c.Pit().Set(e.Car.InPit)
	c.Deslotted().Set(e.Car.Deslotted)
	c.LastLapTime().Set(e.Car.Lap.LapTime)
	c.Laps().Set(uint32(e.Car.Lap.Number))
	c.LastLap().Set(types.Lap{e.Car.Lap.Number, e.Car.Lap.LapTime, e.RaceTimer})
	c.Controller().ButtonTrackCall().Set(e.Car.Controller.TrackCall)
	c.Controller().TriggerValue().Set(uint8(e.Car.Controller.TriggerValue))
}

func (c *Car) Init(ctx context.Context) {
	/*	c.maxSpeed.RegisterObserver(c.maxSpeedChangeObserver)
		c.maxSpeed.Init(ctx)
		c.maxSpeed.Update()

		c.pitLaneMaxSpeed.RegisterObserver(c.pitLaneMaxSpeedChangeObserver)
		c.pitLaneMaxSpeed.Init(ctx)
		c.pitLaneMaxSpeed.Update()

		// c.deslotted.Init(ctx)
		// c.lastLapTime.Init(ctx)
		c.laps.Init(ctx)
		c.lastLap.Init(ctx)

		c.minSpeed.RegisterObserver(func(u uint8, a observable.Annotations) {
			c.implementer.Car(c.id).MinSpeed(u)
		})
		c.minSpeed.Init(ctx)
		c.minSpeed.Update()

		c.maxBreaking.RegisterObserver(c.maxBreakingChangeObserver)
		c.maxBreaking.Init(ctx)
		c.maxBreaking.Update()

		// c.pit.Init(ctx)
		c.controller.Init(ctx)
		c.drivers.Init(ctx)
		c.drivers.Update()*/
}
