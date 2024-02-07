package main

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/devices/oxigen"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)

	device := "/dev/ttyACM0"
	connection, err := oxigen.CreateUSBConnection(&device)
	if err != nil {
		panic(err)
	}

	implement, err := oxigen.CreateImplement(connection)
	if err != nil {
		panic(err)
	}

	err = implement.Start(make(chan<- drivers.Event))
	if err != nil {
		panic(err)
	}

}
