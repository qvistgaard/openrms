package oxigen

import (
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type USBConnection struct {
	device string
}

func NewUSBConnection(device string) *USBConnection {
	o := new(USBConnection)
	o.device = device
	return o
}

func (usb USBConnection) connect() (io.ReadWriteCloser, error) {
	options := serial.OpenOptions{
		PortName:              usb.device,
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	return serial.Open(options)
}
