package pit

/*
import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/state/car"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var pitStopComplete = make(chan bool)

type SimpleTestPitRule struct{}

func (SimpleTestPitRule) HandlePitStop(car *car.Car, cancel <-chan bool) bool {
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
			f := car.Get("testv").(uint8)
			car.Set("testv", f+1)
			pitStopComplete <- true
			return true
		}
	}
}

func setupContext() (*Rule, *state.Car) {
	log.SetLevel(log.DebugLevel)
	ctx := &application.Context{
		Rules: state.CreateRuleList(),
	}
	p := CreatePitRule(ctx)
	ctx.Rules.Append(p)
	ctx.Rules.Append(SimpleTestPitRule{})
	ctx.Course = state.CreateCourse(&state.CourseConfig{}, ctx.Rules)
	ctx.Course.Set(state.RaceStatus, state.RaceStatusRunning)

	car := state.CreateCar(1, nil, ctx.Rules)
	car.Set("testv", uint8(0))

	return p, car
}

func (SimpleTestPitRule) Priority() uint8 {
	return 50
}

func (r SimpleTestPitRule) InitializeCarState(car *state.Car) {
}

func (r SimpleTestPitRule) InitializeCourseState(race *state.Course) {
}

func TestThatPitStopIsStartedWithoutMovingCarAutoConfirmation(t *testing.T) {
	rule, car := setupContext()

	stateMachine := rule.carState[car.Id()]

	car.Set(state.CarInPit, true)
	assert.Equal(t, stateCarInPitLane, stateMachine.MustState())

	select {
	case <-time.After(10 * time.Second):
		log.Info(stateMachine.MustState())
		t.Fatal("pit stop did not complete in time")
	case <-pitStopComplete:
		log.Info("complete channel recieved")
	}
}

func TestThatPitStopIsStartedWithoutMovingCar(t *testing.T) {
	rule, car := setupContext()

	stateMachine := rule.carState[car.Id()]

	car.Set(state.CarInPit, true)
	assert.Equal(t, stateCarInPitLane, stateMachine.MustState())

	car.Set(state.ControllerBtnTrackCall, true)
	assert.Equal(t, stateCarPitStopActive, stateMachine.MustState())
}

func TestThatPitStopIsStartedWithMovingCar(t *testing.T) {
	rule, car := setupContext()

	stateMachine := rule.carState[car.Id()]

	car.Set(state.CarInPit, true)
	assert.Equal(t, stateCarInPitLane, stateMachine.MustState())

	car.Set(state.ControllerTriggerValue, state.TriggerValue(10))
	assert.Equal(t, stateCarMoving, stateMachine.MustState())

	// check that pit stop can't be confirmed while moving
	car.Set(state.ControllerBtnTrackCall, true)
	assert.Equal(t, stateCarMoving, stateMachine.MustState())

	car.Set(state.ControllerBtnTrackCall, false)
	car.Set(state.ControllerTriggerValue, state.TriggerValue(0))
	assert.Equal(t, stateCarStopped, stateMachine.MustState())

	car.Set(state.ControllerBtnTrackCall, true)
	assert.Equal(t, stateCarPitStopActive, stateMachine.MustState())
}

func TestThatPitStopHandling(t *testing.T) {
	rule, car := setupContext()

	stateMachine := rule.carState[car.Id()]

	car.Set(state.CarInPit, true)
	car.Set(state.ControllerBtnTrackCall, true)
	assert.Equal(t, stateCarPitStopActive, stateMachine.MustState())

	select {
	case <-time.After(5 * time.Second):
		t.Fatal("pit stop did not complete in time")
	case <-pitStopComplete:
		log.Info("complete channel recieved")
	}
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, stateCarPitStopComplete, stateMachine.MustState())
}
*/
