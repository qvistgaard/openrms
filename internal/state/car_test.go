package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCarCanBeCreatedAndChangedByReference(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{
		"fuel": 100,
	}, make([]Rule, 0))
	c.State().Get("fuel").Set(80)

	ch := c.State().Changes()
	ch["fuel"].Get()
	assert.Equal(t, 80, ch["fuel"].Get())
}

func TestCarStateWillBeCreatedIfMissing(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{}, make([]Rule, 0))
	c.State().Get("fuel").Set(80)

	ch := c.State().Changes()
	ch["fuel"].Get()
	assert.Equal(t, 80, ch["fuel"].Get())
}

type SimpleRule struct {
}

func (s *SimpleRule) InitializeCarState(car *Car) {
	car.State().Get("test").Set(100)
}

func TestCarWillInitializeRules(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{}, []Rule{
		new(SimpleRule),
	})

	assert.Equal(t, 100, c.State().Get("test").Get())
	assert.True(t, c.State().Get("test").Changed())
}
