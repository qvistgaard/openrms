package serial

import (
	"errors"
	log "github.com/sirupsen/logrus"
	serial "github.com/tarm/goserial"
	"go.bug.st/serial/enumerator"
	"io"
	"strings"
	"time"
)

func CreateUSBConnection(device *string) (io.ReadWriteCloser, error) {
	var oxigenPort string
	if device == nil || *device == "" {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			return nil, errors.New("no serial ports found")
		}
		for _, port := range ports {
			log.WithField("port", port.Name).
				WithField("usb", port.IsUSB).
				WithField("vendor", port.VID).
				WithField("product", port.PID).
				WithField("name", port.Product).
				Debug("found COM port")
			if port.IsUSB && strings.ToUpper(port.VID) == "1FEE" && port.PID == "0002" {
				oxigenPort = port.Name
				log.WithField("port", oxigenPort).
					WithField("vendor", port.VID).
					WithField("product", port.PID).
					Info("oxigen COM port identified")
			}
		}
		if oxigenPort == "" {
			return nil, errors.New("oxigen dongle not found")
		}
	} else {
		oxigenPort = *device
	}
	log.WithField("port", oxigenPort).Info("Using oxigenport")

	options := &serial.Config{Name: oxigenPort, Baud: 115200, ReadTimeout: time.Millisecond * 50}

	// Open the port.
	port, err := serial.OpenPort(options)
	if err != nil {
		return nil, err
	}
	return port, nil
}
