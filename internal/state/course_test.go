package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRaceStateCreation(t *testing.T) {
	race := CreateCourse(&CourseConfig{}, []Rule{
		new(SimpleRule),
	})
	s := race.Get(RaceStatus)

	assert.Equal(t, nil, s)
}
