package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatCommandTypeReturnValidPointerTypeForSwitchLoop(t *testing.T) {
	m := NewMaxSpeed(1, 100)

	assert.Equal(t, uint8(1), m.Driver())
	assert.Equal(t, uint8(100), m.Value()[0])
}
