package v3

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/devices/oxigen/serial"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDriver3xCommunications(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(false)

	connection, err := serial.CreateUSBConnection(nil)
	assert.Nil(t, err)
	assert.NotNil(t, connection)

	implement, err := CreateDriver(connection)
	assert.Nil(t, err)
	assert.NotNil(t, implement)

	received := make(chan drivers.Event)

	implement.Start(received)

	for {
		select {
		case data := <-received:
			log.WithField("data", data).Info("Data received")
		}
	}

}
