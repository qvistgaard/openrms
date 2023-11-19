package car

import (
	"context"
	"github.com/divideandconquer/go-merge/merge"
	"github.com/qvistgaard/openrms/internal/implement"
	config "github.com/qvistgaard/openrms/internal/state/config/car"
	"github.com/qvistgaard/openrms/internal/state/controller"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

func NewCar(implementer implement.Implementer, factory *reactive.Factory, settings *config.CarSettings, defaults *config.CarSettings, id types.Id) *Car {
	a := reactive.Annotations{
		annotations.CarId: id,
	}

	settings = merge.Merge(defaults, settings).(*config.CarSettings)
	car := &Car{
		implementer:     implementer,
		id:              id,
		maxBreaking:     factory.NewDistinctPercent(*settings.MaxBreaking),
		maxSpeed:        factory.NewDistinctPercent(*settings.MaxSpeed, a, reactive.Annotations{annotations.CarValueFieldName: fields.MaxTrackSpeed}),
		minSpeed:        factory.NewDistinctPercent(*settings.MinSpeed, a, reactive.Annotations{annotations.CarValueFieldName: fields.MinSpeed}),
		pitLaneMaxSpeed: factory.NewDistinctPercent(*settings.PitLane.MaxSpeed, a, reactive.Annotations{annotations.CarValueFieldName: fields.MaxPitSpeed}),
		pit:             factory.NewDistinctBoolean(false, a, reactive.Annotations{annotations.CarValueFieldName: fields.InPit}),
		deslotted:       factory.NewDistinctBoolean(false, a, reactive.Annotations{annotations.CarValueFieldName: fields.Deslotted}),
		lastLapTime:     factory.NewDuration(0, a, reactive.Annotations{annotations.CarValueFieldName: fields.LapTime}),
		lastLap:         factory.NewDistinctLapNumber(a, reactive.Annotations{annotations.CarValueFieldName: fields.LastLap}),
		laps:            factory.NewDistinctGauge(0, a, reactive.Annotations{annotations.CarValueFieldName: fields.Laps}),
		drivers:         factory.NewDrivers(*settings.Drivers, a, reactive.Annotations{annotations.CarValueFieldName: fields.Drivers}),
		controller:      controller.NewController(a, factory),
	}
	return car
}

type Car struct {
	id              types.Id
	implementer     implement.Implementer
	controller      *controller.Controller
	pit             *reactive.Boolean
	pitLaneMaxSpeed *reactive.Percent
	maxSpeed        *reactive.Percent
	minSpeed        *reactive.Percent
	maxBreaking     *reactive.Percent
	deslotted       *reactive.Boolean
	lastLapTime     *reactive.Duration
	laps            *reactive.Gauge
	lastLap         *reactive.Lap
	drivers         *reactive.Drivers
}

func (c *Car) PitLaneMaxSpeed() *reactive.Percent {
	return c.pitLaneMaxSpeed
}

func (c *Car) LastLap() *reactive.Lap {
	return c.lastLap
}

func (c *Car) MaxSpeed() *reactive.Percent {
	return c.maxSpeed
}

func (c *Car) MinSpeed() *reactive.Percent {
	return c.minSpeed
}

func (c *Car) Controller() *controller.Controller {
	return c.controller
}

func (c *Car) Id() types.Id {
	return c.id
}

func (c *Car) Pit() *reactive.Boolean {
	return c.pit
}

func (c *Car) Deslotted() *reactive.Boolean {
	return c.deslotted
}

func (c *Car) LastLapTime() *reactive.Duration {
	return c.lastLapTime
}

func (c *Car) Laps() *reactive.Gauge {
	return c.laps
}

func (c *Car) Drivers() *reactive.Drivers {
	return c.drivers
}

func (c *Car) UpdateFromEvent(e implement.Event) {
	c.Pit().Set(e.Car.InPit)
	c.Deslotted().Set(e.Car.Deslotted)
	c.LastLapTime().Set(e.Car.Lap.LapTime)
	c.Laps().Set(float64(e.Car.Lap.Number))
	c.LastLap().Set(types.NewLap(e.Car.Lap.Number, e.RaceTimer))
	c.Controller().ButtonTrackCall().Set(e.Car.Controller.TrackCall)
	c.Controller().TriggerValue().Set(types.NewPercentFromFloat64(e.Car.Controller.TriggerValue))
}

func (c *Car) Init(ctx context.Context, postProcess reactive.ValuePostProcessor) {
	c.maxSpeed.RegisterObserver(c.maxSpeedChangeObserver)
	c.maxSpeed.Init(ctx)
	c.maxSpeed.Update()

	c.pitLaneMaxSpeed.RegisterObserver(c.pitLaneMaxSpeedChangeObserver)
	c.pitLaneMaxSpeed.Init(ctx)
	c.pitLaneMaxSpeed.Update()

	c.deslotted.Init(ctx)
	c.lastLapTime.Init(ctx)
	c.laps.Init(ctx)
	c.lastLap.Init(ctx)

	c.minSpeed.RegisterObserver(c.minSpeedChangeObserver)
	c.minSpeed.Init(ctx)
	c.minSpeed.Update()

	c.maxBreaking.RegisterObserver(c.maxBreakingChangeObserver)
	c.maxBreaking.Init(ctx)
	c.maxBreaking.Update()

	c.pit.Init(ctx)
	c.controller.Init(ctx)
	c.drivers.Init(ctx)
	c.drivers.Update()
}
