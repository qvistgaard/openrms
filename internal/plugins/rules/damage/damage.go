package damage

import (
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

const CarDamage = "damage"

type Damage struct {
}

func (d *Damage) Notify(v *state.Value) {
	if c, ok := v.Owner().(state.Car); ok {
		switch v.Name() {
		case state.CarOnTrack:
			if !v.Get().(bool) {
				d := c.Get(CarDamage).(uint8) + uint8(rand.Int31()) + c.Get(state.ControllerTriggerValue).(uint8)
				c.Set(CarDamage, d)
				if d >= 255 {
					c.Set(limbmode.CarLimbMode, true)
				}
			}
		}
	}
}

func (d *Damage) InitializeRaceState(race *state.Course) {

}

func (d *Damage) InitializeCarState(car *state.Car) {
	m := car.Get(CarDamage)
	if m == nil {
		car.Set(CarDamage, uint8(0))
	}
}

func (d *Damage) HandlePitStop(car *state.Car, cancel chan bool) {
	log.Warn("IMPLEMENT ME")
}

func (d *Damage) Priority() uint8 {
	return 50
}
