package limbmode

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/state"
)

const CarLimbMode = "limb-mode"
const CarLimbModeMaxSpeed = "limb-mode-max-speed"

type Settings struct {
	LimbMode struct {
		MaxSpeed state.MaxSpeed `mapstructure:"max-speed,omitempty"`
	} `mapstructure:"limb-mode"`
}

type LimbMode struct {
	MaxSpeed state.MaxSpeed
}

func (l *LimbMode) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		switch v.Name() {
		case CarLimbMode:
			if v.Get().(bool) {
				c.Set(state.CarMaxSpeed, c.Get(CarLimbModeMaxSpeed))
			} else {
				c.SetDefault(state.CarMaxSpeed)
			}
		case pit.State:
			if v.Get().(string) == pit.Started {
				c.Set(CarLimbMode, false)
			}
		}
	}
}

func (l *LimbMode) InitializeCourseState(race *state.Course) {}

func (l *LimbMode) InitializeCarState(car *state.Car) {
	settings := &Settings{}
	car.Settings(settings)
	m := car.Get(CarLimbMode)
	if m == nil {
		car.Set(CarLimbMode, false)
	}

	ms := car.Get(CarLimbModeMaxSpeed)
	if ms == nil {
		if settings.LimbMode.MaxSpeed > 0 {
			car.Set(CarLimbModeMaxSpeed, settings.LimbMode.MaxSpeed)
		} else {
			car.Set(CarLimbModeMaxSpeed, l.MaxSpeed)
		}
	}
	car.Subscribe(CarLimbMode, l)
	car.Subscribe(pit.State, l)
}
