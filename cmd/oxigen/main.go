package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	serial "github.com/tarm/goserial"
	"io"
	"time"
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)

	device := "/dev/ttyACM0"

	c := &serial.Config{Name: device, Baud: 921600, ReadTimeout: time.Millisecond * 20}
	connection, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	timer := []byte{0x00, 0x00, 0x01}

	if err != nil {
		panic(err)
	}

	for true {

		i := 0
		for i == 0 {
			send := pack(timer)
			write(send, connection)
			_, i = read(connection, timer)
		}
		log.Info("COMPLETE")
		// time.Sleep(100 * time.Millisecond)

	}

	time.Sleep(1 * time.Hour)

}

func read(connection io.ReadWriteCloser, timer []byte) (error, int) {
	buffer := make([]byte, 13)
	r := io.LimitReader(connection, 13)
	n := 0
	var err error

	n, err = r.Read(buffer)
	if err != nil {
		log.Error(err)
		return err, n
	}
	log.WithField("message", fmt.Sprintf("%x", buffer)).
		WithField("bytes", n).
		WithField("time", time.Now()).
		Trace("received message from dongle")

	timer = buffer[9:12]
	return nil, n
}

func write(send []byte, connection io.ReadWriteCloser) error {
	log.WithFields(map[string]interface{}{
		"message": fmt.Sprintf("%x", send),
		"size":    fmt.Sprintf("%d", len(send)),
	}).Trace("send message to dongle")

	_, err := connection.Write(send)
	if err != nil {
		panic(err)
	}
	return err
}

func pack(timer []byte) []byte {
	i := time.Now().UnixMilli() / 10
	log.WithFields(map[string]interface{}{
		"message": fmt.Sprintf("%x", i),
		"t":       fmt.Sprintf("%x", timer),
	}).Trace("timer")

	return []byte{
		0x03 | 0x20 | 0x00,
		0xff,
		0x00,
		0x00,
		0x00,
		0x00,     // unused
		0x00,     // unused
		timer[0], // timer
		timer[1], // timer
		timer[2], // timer
	}

}
