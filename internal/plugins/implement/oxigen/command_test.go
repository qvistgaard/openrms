package oxigen

import (
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPitLaneSpeed(t *testing.T) {
	c := newPitLaneSpeed(1, 200)
	assert.Equal(t, uint8(1), c.id)
	assert.Equal(t, uint8(200), c.value)
}

func TestNewMaxBreaking(t *testing.T) {
	c := newMaxBreaking(1, 200)
	assert.Equal(t, uint8(1), c.id)
	assert.Equal(t, uint8(200), c.value)
}

func TestNewMaxSpeed(t *testing.T) {
	c := newMaxSpeed(1, 200)
	assert.Equal(t, uint8(1), c.id)
	assert.Equal(t, uint8(200), c.value)
}

func TestNewMinSpeedLaneChangeAny(t *testing.T) {
	c := newMinSpeed(1, 200, CarForceLangeChangeAny)

	// Confirm speed divided by 4 when removing
	// lane change bits
	assert.Equal(t, uint8(50), c.value&0x3F)

	// Confirm both lane change bits are set
	assert.Equal(t, uint8(0xC0), c.value&0xC0)

	assert.Equal(t, uint8(1), c.id)
}

func TestNewMinSpeedLaneChangeNone(t *testing.T) {
	c := newMinSpeed(1, 200, CarForceLaneChangeNone)

	// Confirm no lane change bits are set
	assert.Equal(t, uint8(0x00), c.value&0xC0)
	assert.Equal(t, uint8(1), c.id)
}

func TestNewMinSpeedLaneChangeLeft(t *testing.T) {
	c := newMinSpeed(1, 200, CarForceLaneChangeLeft)

	// Confirm no lane change bits are set
	assert.Equal(t, uint8(CarForceLaneChangeLeft), c.value&0xC0)
	assert.Equal(t, uint8(1), c.id)
}

func TestNewMinSpeedLaneChangeRight(t *testing.T) {
	c := newMinSpeed(1, 200, CarForceLaneChangeRight)

	// Confirm no lane change bits are set
	assert.Equal(t, uint8(CarForceLaneChangeRight), c.value&0xC0)
	assert.Equal(t, uint8(1), c.id)
}

func TestRaceCommandSetZeroOnInitialization(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	assert.Equal(t, uint8(0x00), c.state)
}

func TestRaceCommandSetStop(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.stop()
	assert.Equal(t, uint8(0x01), c.state)
}

func TestRaceCommandSetStart(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.start()
	assert.Equal(t, uint8(0x03), c.state)
}

func TestRaceCommandSetPause(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.pause()
	assert.Equal(t, uint8(0x04), c.state)
}

func TestRaceCommandSetFlaggedWithLaneChange(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.flag(true)
	assert.Equal(t, uint8(0x05), c.state)
}

func TestRaceCommandSetFlaggedWithLaneChangeDisabled(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.flag(false)
	assert.Equal(t, uint8(0x15), c.state)
}

func TestRaceCommandSetMaxSpeed(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.maxSpeed(255)
	assert.Equal(t, uint8(0x00), c.state)
	assert.Equal(t, uint8(0xFF), c.settings.maxSpeed)
}

func TestRaceCommandPitLaneLapCountingOnExit(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.pitLaneLapCount(true, false)
	assert.Equal(t, uint8(0x00), c.state)
	assert.Equal(t, uint8(0xFF), c.settings.maxSpeed)
	assert.Equal(t, uint8(0x40), c.settings.pitLane.lapTrigger)
}

func TestRaceCommandPitLaneLapCountingOnEntry(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.pitLaneLapCount(true, true)
	assert.Equal(t, uint8(0x00), c.state)
	assert.Equal(t, uint8(0xFF), c.settings.maxSpeed)
	assert.Equal(t, uint8(0x00), c.settings.pitLane.lapTrigger)
}

func TestRaceCommandPitLaneLapCountingDisabled(t *testing.T) {
	state := state.CourseState{}
	c := newEmptyCommand(state, 0x00, newSettings())
	c.pitLaneLapCount(false, true)
	assert.Equal(t, uint8(0x00), c.state)
	assert.Equal(t, uint8(0xFF), c.settings.maxSpeed)
	assert.Equal(t, uint8(0x20), c.settings.pitLane.lapCounting)
	assert.Equal(t, uint8(0x00), c.settings.pitLane.lapTrigger)
}

func TestCarCommandReturnsFalseIfUnknown(t *testing.T) {
	s := state.CourseState{}
	c := newEmptyCommand(s, 0x00, newSettings())
	v := &state.Value{}
	v.Set(uint8(255))

	b := c.carCommand(1, "unknown", v)
	assert.False(t, b)
}

func createTestValue(n string, v interface{}) (bool, *Command) {
	s := state.CourseState{}
	c := newEmptyCommand(s, 0x00, newSettings())
	b := c.carCommand(1, n, v)

	return b, c
}

func TestCarCommandSetMaxSpeed(t *testing.T) {
	b, c := createTestValue(state.CarMaxSpeed, state.Speed(25))

	assert.True(t, b)
	assert.Equal(t, uint8(25), c.car.value)
}

func TestCarCommandSetMaxBreaking(t *testing.T) {
	b, c := createTestValue(state.CarMaxBreaking, uint8(255))

	assert.True(t, b)
	assert.Equal(t, uint8(255), c.car.value)
}

func TestCarCommandSetMinSpeed(t *testing.T) {
	b, c := createTestValue(state.CarMinSpeed, uint8(255))

	assert.True(t, b)
	assert.Equal(t, uint8(63), c.car.value)
}

func TestCarCommandSetPitLaneSpeed(t *testing.T) {
	b, c := createTestValue(state.CarPitLaneSpeed, state.Speed(255))

	assert.True(t, b)
	assert.Equal(t, uint8(255), c.car.value)
}

func TestRaceStatusChangeFromRaceStateStop(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.RaceStatus, Value: state.RaceStatusStopped},
		},
		Time: time.Now(),
	}
	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x01), c.state)
}

func TestRaceStatusChangeFromRaceStatePaused(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.RaceStatus, Value: state.RaceStatusPaused},
		},
		Time: time.Now(),
	}
	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x04), c.state)
}

func TestRaceStatusChangeFromRaceStateRunning(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.RaceStatus, Value: state.RaceStatusRunning},
		},
		Time: time.Now(),
	}
	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x03), c.state)
}

func TestRaceStatusChangeFromRaceStateFlaggedLCDisabled(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.RaceStatus, Value: state.RaceStatusFlaggedLCDisabled},
		},
		Time: time.Now(),
	}
	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x15), c.state)
}

func TestRaceStatusChangeFromRaceStateFlaggedLCEnabled(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.RaceStatus, Value: state.RaceStatusFlaggedLCEnabled},
		},
		Time: time.Now(),
	}
	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x05), c.state)
}

func TestRaceStatusChangeFromRaceStateMaxSpeed(t *testing.T) {
	s := state.CourseState{
		Changes: []state.Change{
			{Name: state.CourseMaxSpeed, Value: state.Speed(100)},
		},
		Time: time.Now(),
	}

	c := newEmptyCommand(s, 0x00, newSettings())
	assert.Equal(t, uint8(0x64), c.settings.maxSpeed)
}

/*
func TestRaceStateSetToStart(t *testing.T) {
	state := make(map[string]state.StateInterface)
	c := newEmptyCommand(state, 0x00, nil)
	o.Start()
	m := o.command(*ipc.NewEmptyCommand())
	assert.Equal(t, "03000000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToPause(t *testing.T) {
	o := new(Oxigen)
	o.Pause()
	m := o.command(*ipc.NewEmptyCommand())
	assert.Equal(t, "04000000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.Stop()
	o.Speed(255)
	m := o.command(*ipc.NewEmptyCommand())
	assert.Equal(t, "01ff0000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopPitLaneLapCountExitEnabledAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.Stop()
	o.Speed(255)
	o.PitLaneLapCount(true, false)
	m := o.command(*ipc.NewEmptyCommand())

	assert.Equal(t, "41ff0000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopPitLaneLapCountOnEntryEnabledAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.stop()
	o.maxSpeed(255)
	o.pitLaneLapCount(true, true)
	m := o.command(*ipc.NewEmptyCommand())

	assert.Equal(t, "01ff0000000000000000", hex.EncodeToString(m))
}


*/
