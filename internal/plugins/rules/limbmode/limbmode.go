package limbmode

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
)

const CarLimbMode = "limb-mode"
const CarLimbModeMaxSpeed = "limb-mode-max-speed"

type Settings struct {
	LimbMode struct {
		MaxSpeed state.Speed `mapstructure:"max-speed,omitempty"`
	} `mapstructure:"limb-mode"`
}

type LimbMode struct {
	MaxSpeed state.Speed
	course   *state.Course
}

func (l *LimbMode) Notify(v *state.Value) {
	if l.course.Get(state.RaceStatus) != state.RaceStatusStopped {
		if c, ok := v.Owner().(*state.Car); ok {
			switch v.Name() {
			case CarLimbMode:
				if v.Get().(bool) {
					c.Set(state.CarMaxSpeed, c.Get(CarLimbModeMaxSpeed))
					log.WithField("car", c.Id()).
						WithField("speed", c.Get(state.CarMaxSpeed)).
						Debugf("limb-mode: enabled")

				} else {
					c.SetDefault(state.CarMaxSpeed)
					log.WithField("car", c.Id()).
						WithField("speed", c.Get(state.CarMaxSpeed)).
						Debugf("limb-mode: disabled")
				}
			case pit.State:
				if v.Get().(string) == pit.Stopped {
					c.Set(CarLimbMode, false)
				}
			}
		}
	}
}

func (l *LimbMode) InitializeCourseState(course *state.Course) {
	l.course = course
}

func (l *LimbMode) InitializeCarState(car *state.Car) {
	settings := &Settings{}
	car.Settings(settings)
	m := car.Get(CarLimbMode)
	if m == nil {
		car.Set(CarLimbMode, false)
	}

	ms := car.Get(CarLimbModeMaxSpeed)
	if ms == nil {
		// TODO: Fix configuration reading
		if settings.LimbMode.MaxSpeed > 0 {
			car.Set(CarLimbModeMaxSpeed, settings.LimbMode.MaxSpeed)
		} else {
			car.Set(CarLimbModeMaxSpeed, l.MaxSpeed)
		}
	}
	car.Subscribe(CarLimbMode, l)
	car.Subscribe(pit.State, l)
}
