package pit

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/types"
	log "github.com/sirupsen/logrus"
	"testing"
)

type NoopHandler struct {
}

func (n NoopHandler) OnCarStop(stop StartPitStop) error {
	log.Info("OnCarStop called")
	return stop()
}

func (n NoopHandler) Start(complete CompletePitStop, cancel CancelPitStop) error {
	log.Info("Start")
	return complete()
}

func (n NoopHandler) Id() types.CarId {
	return 1
}

func (n NoopHandler) OnCarStart() error {
	log.Info("OnCarStart")
	return nil
}

func (n NoopHandler) OnComplete() error {
	log.Info("OnComplete")
	return nil
}

func (n NoopHandler) Active() observable.Observable[bool] {
	//TODO implement me
	panic("implement me active")
}

func (n NoopHandler) Current() observable.Observable[uint8] {
	//TODO implement me
	panic("implement me current")
}

func Test_machine(t *testing.T) {

	config := Config{}
	config.Plugin.Pit.Enabled = true
	config.Plugin.Pit.Commentary = false

	stateMachine := machine(&NoopHandler{})

	// print(stateMachine.ToGraph())
	var err error

	err = stateMachine.Fire(triggerCarEnteredPitLane)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = stateMachine.Fire(triggerCarExitedPitLane)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = stateMachine.Fire(triggerCarEnteredPitLane)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = stateMachine.Fire(triggerCarMoving)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = stateMachine.Fire(triggerCarStopped)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// log.Info(stateMachine.String())

	// time.Sleep(15 * time.Second)

	// stateMachine.Fire(triggerCarPitStopConfirmed)
	// stateMachine.Fire(triggerCarMoving)

}
