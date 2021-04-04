package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRaceStateCreation(t *testing.T) {
	race := CreateRace(map[string]interface{}{})
	s := race.State().Get(RaceStatus).Get()

	assert.Equal(t, nil, s)
}
