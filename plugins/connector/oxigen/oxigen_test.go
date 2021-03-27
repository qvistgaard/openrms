package oxigen

import (
	"errors"
	"io"
	"log"
	"testing"
)

type MockHandshakeConnection struct {
	input []byte
}

func newMockConnection() *MockHandshakeConnection {
	return new(MockHandshakeConnection)
}

func (mock MockHandshakeConnection) Write(input []byte) (int, error) {
	if input[0] != 0x06 {
		return 0, errors.New("first input byte is not 0x06")
	}
	if input[1] != 0x06 {
		return 0, errors.New("second input byte is not 0x06")
	}
	if input[2] != 0x06 {
		return 0, errors.New("third input byte is not 0x06")
	}
	if input[3] != 0x06 {
		return 0, errors.New("fourth input byte is not 0x06")
	}
	return 13, nil
}
func (mock MockHandshakeConnection) Read(output []byte) (int, error) {
	n := copy(output, []byte{0x03, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return n, nil
}

func (mock MockHandshakeConnection) Close() error {
	return nil
}

func (mock MockHandshakeConnection) connect() (io.ReadWriteCloser, error) {
	return new(MockHandshakeConnection), nil
}

func TestHandshake(t *testing.T) {
	connection := newMockConnection()
	o, err := Connect(*connection)
	if err != nil {
		log.Fatal(err)
	}

	if o.version != "3.10" {
		errors.New("version does not match")
	}

}
