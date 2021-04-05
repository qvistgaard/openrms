package oxigen

import (
	"encoding/hex"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type MockHandshakeConnection struct {
	input  queue.Queue
	output queue.Queue
}

func newMockConnection(input queue.Queue, output queue.Queue) *MockHandshakeConnection {
	c := new(MockHandshakeConnection)
	c.output = output
	c.input = input
	return c
}

func (mock *MockHandshakeConnection) Write(input []byte) (int, error) {
	o := make([]byte, 13)
	n := copy(o, input)
	mock.input.Enqueue(o)
	return n, nil
}
func (mock *MockHandshakeConnection) Read(output []byte) (int, error) {
	if mock.output.GetLen() > 0 {
		i, _ := mock.output.Dequeue()
		n := copy(output, i.([]byte))
		return n, nil
	} else {
		return 0, io.EOF
	}

}

func (mock *MockHandshakeConnection) Close() error {
	return nil
}

func (mock *MockHandshakeConnection) connect() (io.ReadWriteCloser, error) {
	return new(MockHandshakeConnection), nil
}

func TestHandshakeAsksForVersionAndDecodesVersion(t *testing.T) {
	input := queue.NewFIFO()
	output := queue.NewFIFO()
	output.Enqueue([]byte{0x03, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	connection := newMockConnection(input, output)
	o, err := CreateImplement(connection)

	v, _ := input.Dequeue()

	assert.Equal(t, "06060606000000000000000000", hex.EncodeToString(v.([]byte)))
	assert.Nil(t, err)
	assert.EqualValues(t, "3.10", o.version)
}

func TestTestSendSingleCommandOnNoCarStateChanges(t *testing.T) {
	o := new(Oxigen)
	o.settings = newSettings()
	o.commands = queue.NewFIFO()
	race := state.CreateRace(map[string]interface{}{})
	race.State().Get(state.RaceStatus).Set(state.RaceStatusStopped)
	car := state.CreateCar(race, 1, map[string]interface{}{}, make([]state.Rule, 0))

	c := implement.CreateCommand(car)

	o.SendCommand(c)
	assert.Equal(t, 1, o.commands.GetLen())

	dequeue, err := o.commands.Dequeue()
	command := dequeue.(*Command)
	assert.Equal(t, 0, o.commands.GetLen())
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x01), command.state)
}

func TestTestSendSingleCommandOnCarStateChanges(t *testing.T) {
	o := new(Oxigen)
	o.settings = newSettings()
	o.commands = queue.NewFIFO()
	race := state.CreateRace(map[string]interface{}{})
	car := state.CreateCar(race, 1, map[string]interface{}{}, make([]state.Rule, 0))
	car.State().Get(state.CarMaxSpeed).Set(uint8(255))
	c := implement.CreateCommand(car)

	o.SendCommand(c)
	dequeue, err := o.commands.Dequeue()
	command := dequeue.(*Command)
	assert.Equal(t, 0, o.commands.GetLen())
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), command.car.id)
	assert.Equal(t, uint8(0x02), command.car.command)
	assert.Equal(t, uint8(0xFF), command.car.value)
}

func TestEventLoopCanReadMessages(t *testing.T) {
	input := queue.NewFIFO()
	output := queue.NewFIFO()
	c := newEmptyCommand(map[string]state.StateInterface{}, 0x00, newSettings())
	o := Oxigen{
		settings: newSettings(),
		commands: queue.NewFIFO(),
		serial:   newMockConnection(input, output),
	}
	o.commands.Enqueue(c)
	o.EventLoop()
}