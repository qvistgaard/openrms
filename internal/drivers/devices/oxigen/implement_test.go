package oxigen

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
	func TestTestSendSingleCommandOnNoCarStateChanges(t *testing.T) {
		o := new(Oxigen)
		o.settings = newSettings()
		o.commands = make(chan *TX, 10)
		race := state.CreateCourse(&state.CourseConfig{}, &state.RuleList{})
		race.Set(state.RaceStatus, state.RaceStatusStopped)
		car := state.CreateCar(1, map[string]interface{}{}, &state.RuleList{})

		o.SendCarState(car.Changes())

		assert.Equal(t, 0, len(o.commands))
		if len(o.commands) > 0 {
			packCommand := <-o.commands
			assert.Equal(t, 0, len(o.commands))
			assert.Equal(t, uint8(0x01), packCommand.state)
		}
	}

	func TestTestSendSingleCommandOnCarStateChanges(t *testing.T) {
		o := new(Oxigen)
		o.settings = newSettings()
		o.commands = make(chan *TX, 10)
		// race := state.CreateCourse(&state.CourseConfig{}, &state.RuleList{})
		car := state.CreateCar(state.CarId(1), map[string]interface{}{}, &state.RuleList{})
		car.Set(state.CarMaxSpeed, state.Speed(100))

		o.SendCarState(car.Changes())

		packCommand := <-o.commands
		assert.Equal(t, 0, len(o.commands))
		assert.Equal(t, uint8(1), packCommand.car.id)
		assert.Equal(t, uint8(0x82), packCommand.car.packCommand)
		assert.Equal(t, uint8(0x64), packCommand.car.value)
	}

	func TestMaxSpeedCommandWillBeTranslatedToProtocol(t *testing.T) {
		o := new(Oxigen)
		o.settings = newSettings()
		c := &TX{
			settings: Settings{
				maxSpeed: 255,
				pitLane: PitLane{
					lapCounting: 0,
					lapTrigger:  0,
				},
			},
			state: 0,
			car:   newMaxSpeed(1, state.Speed(255)),
		}
		packCommand := o.packCommand(c, []byte{0, 0, 0})

		log.Infof("%s", hex.EncodeToString(packCommand))

}

	func TestEventLoopCanReadMessages(t *testing.T) {
		input := queue.NewFIFO()
		output := queue.NewFIFO()
		c := newEmptyCommand(state.CourseState{}, 0x00, newSettings())
		o := Oxigen{
			settings: newSettings(),
			commands: make(chan *TX, 10),
			serial:   newMockConnection(input, output),
		}
		o.commands <- c
		o.EventLoop()
	}
*/
func TestLapTimeUnpack(t *testing.T) {
	lapTime := unpackLapTime(0, 1)
	assert.Equal(t, int64(10), lapTime.Milliseconds())
}

func TestUnpackRaceTime(t *testing.T) {
	raceTime := unpackRaceTime([4]byte{0, 0, 0, 100}, 0)
	// log.Infof("%s, %f", raceTime.String(), raceTime.Seconds())
	assert.Equal(t, time.Second, raceTime)
}
