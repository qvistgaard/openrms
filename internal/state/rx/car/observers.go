package car

import (
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/reactivex/rxgo/v2"
)

func (c *Car) maxSpeedChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		c.implementer.Car(c.id).MaxSpeed(i.(types.Percent))
	})
}

func (c *Car) minSpeedChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		c.implementer.Car(c.id).MinSpeed(i.(types.Percent))
	})
}

func (c *Car) pitLaneMaxSpeedChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		c.implementer.Car(c.id).PitLaneMaxSpeed(i.(types.Percent))
	})
}
func (c *Car) maxBreakingChangeObserver(observable rxgo.Observable) {
	observable.DoOnNext(func(i interface{}) {
		c.implementer.Car(c.id).MaxBreaking(i.(types.Percent))
	})
}
