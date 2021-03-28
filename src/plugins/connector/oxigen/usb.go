package oxigen

import (
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type USBConnection struct {
	device string
}

func NewUSBConnection(device string) (io.ReadWriteCloser, error) {
	options := serial.OpenOptions{
		PortName:              device,
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 1000,
	}
	return serial.Open(options)
}
