package car

import (
	"context"
	"github.com/qvistgaard/openrms/internal/implement"
	config "github.com/qvistgaard/openrms/internal/state/rx/config/car"
	"github.com/qvistgaard/openrms/internal/state/rx/controller"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/fields"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

func NewCar(implementer implement.Implementer, settings *config.CarSettings, defaults *config.CarSettings, id types.Id) *Car {
	a := reactive.Annotations{
		annotations.CarId: id,
	}

	if settings.MaxSpeed == nil {
		settings.MaxSpeed = defaults.MaxSpeed
	}
	if settings.MaxBreaking == nil {
		settings.MaxBreaking = defaults.MaxBreaking
	}
	if settings.MinSpeed == nil {
		settings.MinSpeed = defaults.MinSpeed
	}
	if settings.PitLane == nil {
		settings.PitLane = defaults.PitLane
	}
	if settings.PitLane.MaxSpeed == nil {
		settings.PitLane.MaxSpeed = defaults.PitLane.MaxSpeed
	}

	car := &Car{
		implementer:     implementer,
		id:              id,
		maxBreaking:     reactive.NewPercent(*settings.MaxBreaking),
		maxSpeed:        reactive.NewPercent(*settings.MaxSpeed, a, reactive.Annotations{annotations.CarValueFieldName: "max-speed"}),
		minSpeed:        reactive.NewPercent(*settings.MinSpeed, a, reactive.Annotations{annotations.CarValueFieldName: "min-speed"}),
		pitLaneMaxSpeed: reactive.NewPercent(*settings.PitLane.MaxSpeed, a, reactive.Annotations{annotations.CarValueFieldName: "pit-lane-max-speed"}),
		pit:             reactive.NewBoolean(false, a, reactive.Annotations{annotations.CarValueFieldName: fields.InPit}),
		deslotted:       reactive.NewBoolean(false, a, reactive.Annotations{annotations.CarValueFieldName: fields.Deslotted}),
		lastLapTime:     reactive.NewDuration(0, a, reactive.Annotations{annotations.CarValueFieldName: fields.LapTime}),
		laps:            reactive.NewGauge(0, a, reactive.Annotations{annotations.CarValueFieldName: fields.Laps}),
		controller:      controller.NewController(a),
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
}

func (c *Car) MaxSpeed() *reactive.Percent {
	return c.maxSpeed
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

func (c *Car) UpdateFromEvent(e implement.Event) {
	c.Pit().Set(e.Car.InPit)
	c.Deslotted().Set(e.Car.Deslotted)
	c.LastLapTime().Set(e.Car.Lap.LapTime)
	c.Laps().Set(float64(e.Car.Lap.Number))
	c.Controller().ButtonTrackCall().Set(e.Car.Controller.TrackCall)
	c.Controller().TriggerValue().Set(types.NewPercentFromFloat64(e.Car.Controller.TriggerValue))
}

func (c *Car) Init(ctx context.Context, postProcess reactive.ValuePostProcessor) {
	c.maxSpeed.RegisterObserver(c.maxSpeedChangeObserver)
	c.maxSpeed.Init(ctx, postProcess)
	c.maxSpeed.Update()

	c.pitLaneMaxSpeed.RegisterObserver(c.pitLaneMaxSpeedChangeObserver)
	c.pitLaneMaxSpeed.Init(ctx, postProcess)
	c.pitLaneMaxSpeed.Update()

	c.deslotted.Init(ctx, postProcess)
	c.lastLapTime.Init(ctx, postProcess)
	c.laps.Init(ctx, postProcess)

	c.minSpeed.RegisterObserver(c.minSpeedChangeObserver)
	c.minSpeed.Init(ctx, postProcess)
	c.minSpeed.Update()

	c.maxBreaking.RegisterObserver(c.maxBreakingChangeObserver)
	c.maxBreaking.Init(ctx, postProcess)
	c.maxBreaking.Update()

	c.pit.Init(ctx, postProcess)
	c.controller.Init(ctx, postProcess)
}
