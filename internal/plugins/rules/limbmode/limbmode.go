package limbmode

import "github.com/qvistgaard/openrms/internal/state"

const CarLimbMode = "limb-mode"
const CarLimbModeMaxSpeed = "limb-mode-max-speed"

type LimbMode struct {
}

func (l *LimbMode) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		switch v.Name() {
		case CarLimbMode:
			if v.Get().(bool) {
				c.Set(state.CarMaxSpeed, c.Get(CarLimbModeMaxSpeed))
			} else {
				c.SetDefault(state.CarMaxSpeed)
			}
		}
	}
}

func (l *LimbMode) InitializeRaceState(race *state.Course) {

}

func (l *LimbMode) InitializeCarState(car *state.Car) {
	m := car.Get(CarLimbMode)
	if m == nil {
		car.Set(CarLimbMode, false)
	}

	ms := car.Get(CarLimbModeMaxSpeed)
	if ms == nil {
		car.Set(CarLimbModeMaxSpeed, 100)
	}
	car.Subscribe(CarLimbMode, l)
}
