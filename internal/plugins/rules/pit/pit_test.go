package pit

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

type SimpleTestPitRule struct{}

func (SimpleTestPitRule) HandlePitStop(car *state.Car, cancel <-chan bool) bool {
	log.Info("Ruinning test handler")
	for {
		select {
		case v := <-cancel:
			log.WithField("car", car.Id()).
				WithField("cancel", v).
				WithField("length", len(cancel)).
				Infof("CENCELLED")
			return false
		case <-time.After(500 * time.Millisecond):
			log.Info("counint")
			f := car.Get("testv").(uint8)
			car.Set("testv", f+1)
			log.Infof("Run PIT %+v", f)
			return true
		}
	}
}

func (SimpleTestPitRule) Priority() uint8 {
	return 50
}

func (r SimpleTestPitRule) InitializeCarState(car *state.Car) {
}

func (r SimpleTestPitRule) InitializeCourseState(race *state.Course) {
}

func TestSomething(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := &context.Context{
		Rules: state.CreateRuleList(),
	}
	p := CreatePitRule(ctx)
	ctx.Rules.Append(p)
	ctx.Rules.Append(SimpleTestPitRule{})
	ctx.Course = state.CreateCourse(&state.CourseConfig{}, ctx.Rules)
	ctx.Course.Set(state.RaceStatus, state.RaceStatusRunning)

	c := state.CreateCar(1, nil, ctx.Rules)
	c.Set(state.CarInPit, true)
	c.Set("testv", uint8(0))
	log.Info("1 --------------------------------------------")
	c.Set(state.ControllerTriggerValue, state.TriggerValue(0))
	time.Sleep(1 * time.Second)
	log.Info("2 --------------------------------------------")

	c.Set(state.ControllerTriggerValue, state.TriggerValue(5))
	c.Set(state.CarInPit, false)
	log.Info("3 --------------------------------------------")

	c.Set(state.CarInPit, true)
	c.Set("testv", uint8(0))

	c.Set(state.ControllerTriggerValue, state.TriggerValue(0))
	time.Sleep(1 * time.Second)
	log.Info("4 --------------------------------------------")

	c.Set(state.ControllerTriggerValue, state.TriggerValue(5))
	c.Set(state.CarInPit, false)
	log.Info("5 --------------------------------------------")

}

func TestChannel(t *testing.T) {
	c := make(chan bool, 1)
	close(c)
	c <- true

}
