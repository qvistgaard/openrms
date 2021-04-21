package oxigen

import (
	"encoding/hex"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
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
	car := state.CreateCar(race, 1, map[string]interface{}{}, &state.RuleList{})

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
	race := state.CreateCourse(&state.CourseConfig{}, &state.RuleList{})
	car := state.CreateCar(race, 1, map[string]interface{}{}, &state.RuleList{})
	car.Set(state.CarMaxSpeed, uint8(255))

	o.SendCarState(car.Changes())

	command := <-o.commands
	assert.Equal(t, 0, len(o.commands))
	assert.Equal(t, uint8(1), command.car.id)
	assert.Equal(t, uint8(0x02), command.car.command)
	assert.Equal(t, uint8(0xFF), command.car.value)
}

func TestEventLoopCanReadMessages(t *testing.T) {
	input := queue.NewFIFO()
	output := queue.NewFIFO()
	c := newEmptyCommand(state.CourseChanges{}, 0x00, newSettings())
	o := Oxigen{
		settings: newSettings(),
		commands: make(chan *Command, 10),
		serial:   newMockConnection(input, output),
	}
	o.commands <- c
	o.EventLoop()
}

func TestRaceTimerTranslation(t *testing.T) {
	b := []byte{0x00, 0x00, 0xe8, 0x0b, 0x09}
	r := ((uint(b[0]) * 16777216) + (uint(b[1]) * 65536) + (uint(b[2]) * 256) + uint(b[3])) - uint(b[4])
	ms := time.Duration(r * 10)
	log.Infof("%d", r)
	log.Infof("%d", ms)
	log.Infof("%s", ms*time.Millisecond)
	log.Infof("%s", 593940*time.Millisecond)
	// TODO: FIX RACE TIMER

}
