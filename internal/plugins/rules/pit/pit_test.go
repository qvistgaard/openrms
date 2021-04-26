package pit

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

type SimpleTestPitRule struct{}

func (SimpleTestPitRule) HandlePitStop(car *state.Car, cancel chan bool) {
	log.Info("Ruinning test handler")
	for {
		select {
		case <-cancel:
			log.Info("cancelled")
			return
		case <-time.After(500 * time.Millisecond):
			log.Info("counint")
			f := car.Get("testv").(uint8)
			car.Set("testv", f+1)
			log.Infof("Run PIT %+v", f)
		}
	}
}

func (SimpleTestPitRule) Priority() uint8 {
	// panic("implement me")
	return 50
}

func (r SimpleTestPitRule) InitializeCarState(car *state.Car) {
}

func (r SimpleTestPitRule) InitializeCourseState(race *state.Course) {
}

func TestSomething(t *testing.T) {

	ctx := &context.Context{
		Rules: state.CreateRuleList(),
	}
	p := CreatePitRule(ctx)
	ctx.Rules.Append(p)
	ctx.Rules.Append(SimpleTestPitRule{})

	c := state.CreateCar(1, nil, ctx.Rules)
	c.Set(state.CarInPit, true)
	c.Set("testv", uint8(0))

	c.Set(state.ControllerTriggerValue, state.TriggerValue(0))

	time.Sleep(5 * time.Second)
	c.Set(state.ControllerTriggerValue, state.TriggerValue(5))

}

func TestChannel(t *testing.T) {
	c := make(chan bool, 1)
	close(c)
	c <- true

}
