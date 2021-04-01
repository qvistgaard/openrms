package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type StateChangeConfirmer struct {
}

func TestValue_InitalStateNotChanged(t *testing.T) {
	s := createState(nil, "fuel", 100)
	assert.Equal(t, 100, s.Get())
	assert.False(t, s.Changed())
}

func TestUninitializedValue_StateChanged(t *testing.T) {
	s := createState(nil, "fuel", 100)
	s.Set(101)

	assert.Equal(t, 101, s.Get())
	assert.False(t, s.Changed())
}

func TestValue_StateChanged(t *testing.T) {
	s := createState(nil, "fuel", 100)
	s.initialize()
	s.Set(101)

	assert.Equal(t, 101, s.Get())
	assert.True(t, s.Changed())
}

func TestValue_StateChangedReset(t *testing.T) {
	s := createState(nil, "fuel", 100)
	s.Set(101)

	s.reset()
	assert.False(t, s.Changed())

}

func TestValue_StateInitialIsNotChanged(t *testing.T) {
	s := createState(nil, "fuel", 100)
	s.Set(101)

	assert.Equal(t, 100, s.Initial())
}
