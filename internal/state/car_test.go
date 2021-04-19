package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCarCanBeCreatedAndChangedByReference(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{
		"fuel": 100,
	}, make([]Rule, 0))
	c.Set("fuel", 80)

	ch := c.Changes()
	assert.Equal(t, 80, ch.Changes[0].Value)
	assert.Equal(t, "fuel", ch.Changes[0].Name)
}

func TestCarStateWillBeCreatedIfMissing(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{}, make([]Rule, 0))
	c.Set("fuel", 80)

	ch := c.Changes()
	assert.Equal(t, 80, ch.Changes[0].Value)
	assert.Equal(t, "fuel", ch.Changes[0].Name)
}

type SimpleRule struct {
}

func (s *SimpleRule) InitializeCarState(car *Car) {
	car.Set("test", 100)
}

func (s *SimpleRule) InitializeRaceState(race *Course) {

}

func TestCarWillInitializeRules(t *testing.T) {
	c := CreateCar(nil, 1, map[string]interface{}{}, []Rule{
		new(SimpleRule),
	})

	ch := c.Changes()
	assert.Equal(t, 100, c.Get("test"))
	assert.Equal(t, "test", ch.Changes[0].Name)
}
