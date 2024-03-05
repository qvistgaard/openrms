package serial

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"strings"
)

func CreateUSBConnection(device *string) (serial.Port, error) {
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

	/*	options := &serial.Config{Name: oxigenPort, Baud: 115200, ReadTimeout: time.Millisecond * 50}
	 */
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	// Open the port.
	port, err := serial.Open(oxigenPort, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}
