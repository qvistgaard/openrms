package leaderboard

import (
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatLastLapsWillUpdateCorrectly(t *testing.T) {
	l0 := &LastLapDefault{
		Laps: make([]state.Lap, 2),
	}

	l1 := l0.update(state.Lap{
		LapNumber: 1,
		RaceTimer: state.RaceTimer(1),
		LapTime:   state.LapTime(1),
	})

	l2 := l1.update(state.Lap{
		LapNumber: 2,
		RaceTimer: state.RaceTimer(2),
		LapTime:   state.LapTime(2),
	})

	l3 := l2.update(state.Lap{
		LapNumber: 3,
		RaceTimer: state.RaceTimer(3),
		LapTime:   state.LapTime(3),
	})

	assert.True(t, l0 != l1)
	assert.True(t, l1 != l2)
	assert.True(t, l2 != l3)
	assert.NotEqual(t, l0, l1)
	assert.NotEqual(t, l1, l2)
	assert.NotEqual(t, l2, l3)

	assert.Equal(t, state.LapNumber(0), l0.Laps[0].LapNumber)
	assert.Equal(t, state.RaceTimer(0), l0.Laps[0].RaceTimer)
	assert.Equal(t, state.LapTime(0), l0.Laps[0].LapTime)
	assert.Equal(t, state.LapNumber(0), l0.Laps[1].LapNumber)
	assert.Equal(t, state.RaceTimer(0), l0.Laps[1].RaceTimer)
	assert.Equal(t, state.LapTime(0), l0.Laps[1].LapTime)
	assert.Equal(t, 2, len(l0.Laps))

	assert.Equal(t, state.LapNumber(1), l1.Laps[0].LapNumber)
	assert.Equal(t, state.RaceTimer(1), l1.Laps[0].RaceTimer)
	assert.Equal(t, state.LapTime(1), l1.Laps[0].LapTime)
	assert.Equal(t, state.LapNumber(0), l1.Laps[1].LapNumber)
	assert.Equal(t, state.RaceTimer(0), l1.Laps[1].RaceTimer)
	assert.Equal(t, state.LapTime(0), l1.Laps[1].LapTime)
	assert.Equal(t, 2, len(l1.Laps))

	assert.Equal(t, state.LapNumber(2), l2.Laps[0].LapNumber)
	assert.Equal(t, state.RaceTimer(2), l2.Laps[0].RaceTimer)
	assert.Equal(t, state.LapTime(2), l2.Laps[0].LapTime)
	assert.Equal(t, state.LapNumber(1), l2.Laps[1].LapNumber)
	assert.Equal(t, state.RaceTimer(1), l2.Laps[1].RaceTimer)
	assert.Equal(t, state.LapTime(1), l2.Laps[1].LapTime)
	assert.Equal(t, 2, len(l2.Laps))

	assert.Equal(t, state.LapNumber(3), l3.Laps[0].LapNumber)
	assert.Equal(t, state.RaceTimer(3), l3.Laps[0].RaceTimer)
	assert.Equal(t, state.LapTime(3), l3.Laps[0].LapTime)
	assert.Equal(t, state.LapNumber(2), l3.Laps[1].LapNumber)
	assert.Equal(t, state.RaceTimer(2), l3.Laps[1].RaceTimer)
	assert.Equal(t, state.LapTime(2), l3.Laps[1].LapTime)
	assert.Equal(t, 2, len(l3.Laps))
}

func TestLastLapDefaultCanCompare(t *testing.T) {
	l0 := &LastLapDefault{
		Laps: make([]state.Lap, 2),
	}

	l1 := &LastLapDefault{
		Laps: make([]state.Lap, 2),
	}
	l2 := &LastLapDefault{
		Laps: []state.Lap{
			{LapNumber: state.LapNumber(1), RaceTimer: state.RaceTimer(1), LapTime: state.LapTime(1)},
		},
	}

	assert.True(t, l0.Compare(l1))
	assert.False(t, l0.Compare(l2))
}

func TestCarInitialization(t *testing.T) {
	c := state.CreateCar(state.CarId(1), nil, nil, &state.RuleList{})
	b := &Rule{}
	b.InitializeCarState(c)

	lapDefault, ok := c.Get(CarLastLaps).(*LastLapDefault)
	assert.True(t, ok)
	assert.Equal(t, state.LapNumber(0), lapDefault.Laps[0].LapNumber)
}

func TestRaceInitialization(t *testing.T) {
	config := &state.CourseConfig{}
	c := state.CreateCourse(config, &state.RuleList{})
	b := &Rule{}
	b.InitializeCourseState(c)

	board, ok := c.Get(RaceLeaderboard).(*Default)
	assert.Equal(t, 0, len(board.Entries))
	assert.True(t, ok)
}

func TestCarLapNotification(t *testing.T) {
	config := &state.CourseConfig{}
	co := state.CreateCourse(config, &state.RuleList{})
	c := state.CreateCar(state.CarId(1), nil, nil, &state.RuleList{})
	v := state.CreateState(c, state.CarLap, state.Lap{
		LapNumber: 1,
		RaceTimer: 2,
		LapTime:   3,
	})

	b := &Rule{}
	b.InitializeCarState(c)
	b.InitializeCourseState(co)
	b.Notify(v)

	laps := c.Get(CarLastLaps).(*LastLapDefault)
	lb := co.Get(RaceLeaderboard).(*Default)
	pos := c.Get(CarPosition).(Position)
	assert.Equal(t, state.LapNumber(1), laps.Laps[0].LapNumber)
	assert.Equal(t, state.LapNumber(1), lb.Entries[0].Lap.LapNumber)
	assert.Equal(t, Position(1), pos)

}

func TestLeaderBoardDefaultSorting(t *testing.T) {
	l0 := &Default{Entries: []BoardEntry{}}
	l1, _ := l0.updateCar(state.CarId(1), state.Lap{
		LapNumber: 0, RaceTimer: 2, LapTime: 0,
	})
	l2, _ := l1.updateCar(state.CarId(2), state.Lap{
		LapNumber: 1, RaceTimer: 3, LapTime: 0,
	})
	l3, _ := l2.updateCar(state.CarId(3), state.Lap{
		LapNumber: 1, RaceTimer: 2, LapTime: 0,
	})

	c := l3.(*Default)
	assert.Equal(t, state.CarId(3), c.Entries[0].Car)
	assert.Equal(t, state.CarId(2), c.Entries[1].Car)
	assert.Equal(t, state.CarId(1), c.Entries[2].Car)
}
