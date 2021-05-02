package tirewear

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	"math"
)

const CarTireWear = "tire-wear"

type TireWear float32

type Rule struct {
}

func (t *Rule) Initialize() {
	// panic("implement me")
}

func (t *Rule) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		switch v.Name() {
		case state.ControllerTriggerValue:
			cv := float64(v.Get().(state.TriggerValue))
			ot := c.Get(state.CarOnTrack).(bool)
			ip := c.Get(state.CarInPit).(bool)
			tw := float64(c.Get(CarTireWear).(TireWear))
			if ot && !ip {
				if cv == 0 { // Assume breaking when not in pit
					tw = tw + 1
				}
				if cv > 0 {
					tw = tw + (math.Sqrt(cv) / 1000)
				}
				c.Set(CarTireWear, TireWear(tw))
				if tw >= 100 {
					c.Set(limbmode.CarLimbMode, true)
				}
			}
		}
	}
}

func (t *Rule) InitializeCourseState(race *state.Course) {

}

func (t *Rule) HandlePitStop(car *state.Car, cancel <-chan bool) bool {
	car.Set(CarTireWear, TireWear(0))
	return true
}

func (t *Rule) Priority() uint8 {
	return 75
}

func (t *Rule) InitializeCarState(car *state.Car) {
	m := car.Get(CarTireWear)
	if m == nil {
		car.Set(CarTireWear, TireWear(0))
	}
}
