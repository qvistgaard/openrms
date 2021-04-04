package implement

import (
	"github.com/stretchr/testify/assert"
	"openrms/state"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	c := state.CreateCar(state.CreateRace(nil), 1, nil, nil)
	c.State().Get("test").Set(100)

	command := CreateCommand(c)
	s := command.Changes.Car["test"].Get()
	assert.Equal(t, 100, s)
}
