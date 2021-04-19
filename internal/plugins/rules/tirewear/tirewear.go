package tirewear

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	"math"
)

const CarTireWear = "tire-wear"

type TireWear struct {
}

func (t *TireWear) Initialize() {
	panic("implement me")
}

func (t *TireWear) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		switch v.Name() {
		case state.ControllerTriggerValue:
			cv := float64(v.Get().(uint8))
			ot := c.Get(state.CarOnTrack).(bool)
			ip := c.Get(state.CarInPit).(bool)
			tw := float64(c.Get(CarTireWear).(float32))
			if ot && !ip {
				if cv == 0 { // Assume breaking when not in pit
					tw = tw + 1
				}
				if cv > 0 {
					tw = tw + (math.Sqrt(cv) / 1000)
				}
				c.Set(CarTireWear, float32(tw))
				if tw >= 100 {
					c.Set(limbmode.CarLimbMode, true)
				}
			}
		}
	}
}

func (t *TireWear) InitializeRaceState(race *state.Course) {

}

func (t *TireWear) InitializeCarState(car *state.Car) {
	m := car.Get(CarTireWear)
	if m == nil {
		car.Set(CarTireWear, float32(0))
	}
}
