package state

import (
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/stretchr/testify/assert"
	"openrms/ipc"
	"openrms/plugins/connector"
	"testing"
)

func TestCarCanBeCreatedAndChangedByReference(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{
		"fuel": 100,
	}, make([]Rule, 0))
	c.Get("fuel").Set(80)

	ch := c.StateChanges()
	ch["fuel"].Get()
	assert.Equal(t, 80, ch["fuel"].Get())
}

func TestCarStateWillBeCreatedIfMissing(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{}, make([]Rule, 0))
	c.Get("fuel").Set(80)

	ch := c.StateChanges()
	ch["fuel"].Get()
	assert.Equal(t, 80, ch["fuel"].Get())
}

type SimpleRule struct {
}

func (s *SimpleRule) InitializeCarState(car *Car) {
	car.Get("test").Set(100)
}

func (s *SimpleRule) Handle(race connector.Connector, telemetry queue.Queue, car *Car, event ipc.Event) {

}

func TestCarWillInitializeRules(t *testing.T) {
	c := CreateCar(1, map[string]interface{}{}, []Rule{
		new(SimpleRule),
	})

	assert.Equal(t, 100, c.Get("test").Get())
	assert.True(t, c.Get("test").Changed())
}
