package car

import (
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestCarCanBeCreatedAndChangedByReference(t *testing.T) {

	car := NewCar(nil, nil, nil)
	car.Pit().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			log.Infof("i: %+v", i)
		})
	})
	car.Pit().RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			log.Infof("b: %+v", i)
		})
	})

	car.Pit().Init(nil)

	car.Pit().Set(true)
	car.Pit().Set(false)
	car.Pit().Set(true)

	time.Sleep(5 * time.Second)
}
