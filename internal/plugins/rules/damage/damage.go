package damage

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
)

const CarDamage = "damage"

type Damage uint8

type Rule struct {
}

func (r *Rule) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		switch v.Name() {
		case state.CarOnTrack:
			if ot, ok := v.Get().(bool); ok && !ot {
				t := c.Get(state.ControllerTriggerValue).(state.TriggerValue)
				d := c.Get(CarDamage).(Damage)
				nd := uint8(d) + uint8(t)
				c.Set(CarDamage, Damage(nd))
				if d >= 255 {
					c.Set(limbmode.CarLimbMode, true)
				}
			}
		}
	}
}

func (r *Rule) InitializeCourseState(race *state.Course) {

}

func (r *Rule) InitializeCarState(car *state.Car) {
	m := car.Get(CarDamage)
	if m == nil {
		car.Set(CarDamage, Damage(0))
	}
}

func (r *Rule) HandlePitStop(car *state.Car, cancel <-chan bool) bool {
	// TODO: Handle pit stop correctly locking the car for a certain duration of time
	car.Set(CarDamage, Damage(0))
	return true
}

func (r *Rule) Priority() uint8 {
	return 50
}
