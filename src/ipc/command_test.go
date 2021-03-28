package ipc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCommandType struct {
}

func (t *TestCommandType) Value() []byte {
	return []byte{uint8(1)}
}

func TestThatCommandTypeReturnValidPointerTypeForSwitchLoop(t *testing.T) {
	c := NewCommand(1, new(TestCommandType))
	var swCheck bool = false

	assert.Equal(t, uint8(1), c.Value()[0])
	assert.Equal(t, uint8(1), c.Driver())
	switch c.CommandType().(type) {
	case *TestCommandType:
		swCheck = true
	}
	assert.True(t, swCheck)

}
