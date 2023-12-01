package pit

import (
	"testing"
	"time"
)

func Test_machine(t *testing.T) {

	stateMachine := machine(1, &DefaultHandler{
		cancel: make(chan bool, 1),
	})

	stateMachine.Fire(triggerCarEnteredPitLane)
	stateMachine.Fire(triggerCarExitedPitLane)

	stateMachine.Fire(triggerCarEnteredPitLane)
	stateMachine.Fire(triggerCarMoving)
	stateMachine.Fire(triggerCarStopped)
	time.Sleep(15 * time.Second)

	// stateMachine.Fire(triggerCarPitStopConfirmed)
	// stateMachine.Fire(triggerCarMoving)

}
