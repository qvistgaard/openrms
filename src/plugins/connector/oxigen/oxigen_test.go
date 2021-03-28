package oxigen

import (
	"../../../ipc"
	"encoding/hex"
	queue "github.com/enriquebris/goconcurrentqueue"
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
	o, err := Connect(connection)

	v, _ := input.Dequeue()

	assert.Equal(t, "06060606000000000000000000", hex.EncodeToString(v.([]byte)))
	assert.Nil(t, err)
	assert.EqualValues(t, "3.10", o.version)
}

func TestRaceStateSetToStop(t *testing.T) {
	o := new(Oxigen)
	o.stop()
	m := o.message(*ipc.NewEmptyCommand())
	assert.Equal(t, "01000000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStart(t *testing.T) {
	o := new(Oxigen)
	o.start()
	m := o.message(*ipc.NewEmptyCommand())
	assert.Equal(t, "03000000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToPause(t *testing.T) {
	o := new(Oxigen)
	o.pause()
	m := o.message(*ipc.NewEmptyCommand())
	assert.Equal(t, "04000000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.stop()
	o.maxSpeed(255)
	m := o.message(*ipc.NewEmptyCommand())
	assert.Equal(t, "01ff0000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopPitLaneLapCountExitEnabledAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.stop()
	o.maxSpeed(255)
	o.pitLaneLapCount(true, false)
	m := o.message(*ipc.NewEmptyCommand())

	assert.Equal(t, "41ff0000000000000000", hex.EncodeToString(m))
}

func TestRaceStateSetToStopPitLaneLapCountOnEntryEnabledAndMaxSpeedFull(t *testing.T) {
	o := new(Oxigen)
	o.stop()
	o.maxSpeed(255)
	o.pitLaneLapCount(true, true)
	m := o.message(*ipc.NewEmptyCommand())

	assert.Equal(t, "01ff0000000000000000", hex.EncodeToString(m))
}

func TestEventLoopPublishingMessagesOnQueue(t *testing.T) {
	mockInput := queue.NewFIFO()
	mockOutput := queue.NewFIFO()
	mockOutput.Enqueue([]byte{0x03, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	o := &Oxigen{
		serial: newMockConnection(mockInput, mockOutput),
	}
	o.stop()

	eventInput := queue.NewFIFO()
	eventInput.Enqueue(ipc.NewEmptyCommand())
	eventOutput := queue.NewFIFO()

	o.EventLoop(eventInput, eventOutput)

	var err error

	mockInputValue, err := mockInput.Dequeue()
	assert.Equal(t, "01000000000000000000000000", hex.EncodeToString(mockInputValue.([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 0, mockInput.GetLen()) // confirm no other messages was sent to serial mock

	// TODO: Change this once Event has been implemented
	eventOutputValue, err := eventOutput.Dequeue()
	assert.Equal(t, "030a0000000000000000000000", hex.EncodeToString(eventOutputValue.([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 0, eventOutput.GetLen()) // confirm no other messages was sent to serial mock

	assert.Equal(t, 0, eventInput.GetLen())
}
