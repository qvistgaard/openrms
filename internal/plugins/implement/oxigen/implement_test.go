package oxigen

import (
	"encoding/hex"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
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
	o.commands = make(chan *Command, 10)
	race := state.CreateCourse(&state.CourseConfig{}, &state.RuleList{})
	race.Set(state.RaceStatus, state.RaceStatusStopped)
	car := state.CreateCar(1, map[string]interface{}{}, nil, &state.RuleList{})

	o.SendCarState(car.Changes())

	assert.Equal(t, 0, len(o.commands))
	if len(o.commands) > 0 {
		command := <-o.commands
		assert.Equal(t, 0, len(o.commands))
		assert.Equal(t, uint8(0x01), command.state)
	}
}

func TestTestSendSingleCommandOnCarStateChanges(t *testing.T) {
	o := new(Oxigen)
	o.settings = newSettings()
	o.commands = make(chan *Command, 10)
	// race := state.CreateCourse(&state.CourseConfig{}, &state.RuleList{})
	car := state.CreateCar(state.CarId(1), map[string]interface{}{}, nil, &state.RuleList{})
	car.Set(state.CarMaxSpeed, state.Speed(100))

	o.SendCarState(car.Changes())

	command := <-o.commands
	assert.Equal(t, 0, len(o.commands))
	assert.Equal(t, uint8(1), command.car.id)
	assert.Equal(t, uint8(0x82), command.car.command)
	assert.Equal(t, uint8(0x64), command.car.value)
}

func TestMaxSpeedCommandWillBeTranslatedToProtocol(t *testing.T) {
	o := new(Oxigen)
	o.settings = newSettings()
	c := &Command{
		settings: Settings{
			maxSpeed: 255,
			pitLane: PitLane{
				lapCounting: 0,
				lapTrigger:  0,
			},
		},
		state: 0,
		car:   newMaxSpeed(1, state.Speed(255)),
	}
	command := o.command(c, []byte{0, 0, 0})

	log.Infof("%s", hex.EncodeToString(command))

}

func TestEventLoopCanReadMessages(t *testing.T) {
	input := queue.NewFIFO()
	output := queue.NewFIFO()
	c := newEmptyCommand(state.CourseState{}, 0x00, newSettings())
	o := Oxigen{
		settings: newSettings(),
		commands: make(chan *Command, 10),
		serial:   newMockConnection(input, output),
	}
	o.commands <- c
	o.EventLoop()
}
