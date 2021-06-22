package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRaceStateCreation(t *testing.T) {
	race := CreateCourse(&CourseConfig{}, &RuleList{
		rules:    []Rule{new(SimpleRule)},
		pitRules: nil,
	})
	s := race.Get(RaceStatus)

	assert.Equal(t, uint8(0), s)
}
