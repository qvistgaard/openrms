package state

/*
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCarCanBeCreatedAndChangedByReference(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{
		"fuel": 100,
	}, &RuleList{})
	c.Set("fuel", 80)

	ch := c.Changes()
	assert.Equal(t, 80, ch.Changes[0].Value)
	assert.Equal(t, "fuel", ch.Changes[0].Name)
}

func TestCarStateWillBeCreatedIfMissing(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{}, &RuleList{})
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

func (s *SimpleRule) InitializeCourseState(race *Course) {

}

func TestCarWillInitializeRules(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{}, &RuleList{
		rules:    []Rule{new(SimpleRule)},
		pitRules: nil,
	})

	ch := c.Changes()
	assert.Equal(t, 100, c.Get("test"))
	assert.Equal(t, "test", ch.Changes[0].Name)
}
*/
