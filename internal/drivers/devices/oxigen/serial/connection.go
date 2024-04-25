package serial

import (
	"errors"
	"github.com/rs/zerolog"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"strings"
)

func CreateUSBConnection(logger zerolog.Logger, device *string) (serial.Port, error) {
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
			logger.Debug().
				Str("port", port.Name).
				Bool("usb", port.IsUSB).
				Str("vendor", port.VID).
				Str("product", port.PID).
				Str("name", port.Product).
				Msg("found COM port")
			if port.IsUSB && strings.ToUpper(port.VID) == "1FEE" && port.PID == "0002" {
				oxigenPort = port.Name
				logger.Info().
					Str("port", oxigenPort).
					Str("vendor", port.VID).
					Str("product", port.PID).
					Msg("oxigen COM port identified")
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
		BaudRate: 9600,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		DataBits: 8,
	}

	// Open the port.
	port, err := serial.Open(oxigenPort, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}
