package main

import (
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	serial "github.com/tarm/goserial"
	"go.bug.st/serial/enumerator"
	"io"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)

	f, err := os.Create("profile.gproff")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()
	defer f.Close() // error handling omitted for example

	var oxigenPort string

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		panic(err)
	}
	for _, port := range ports {
		if port.IsUSB && strings.ToUpper(port.VID) == "1FEE" && port.PID == "0002" {
			oxigenPort = port.Name
		}
	}
	if oxigenPort == "" {
		panic(err)
	}

	c := &serial.Config{Name: oxigenPort, Baud: 115200, ReadTimeout: time.Millisecond * 50}
	connection, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	timer := []byte{0x00, 0x00, 0x01}

	if err != nil {
		panic(err)
	}

	links := make(map[uint32]link)
	expire := make(chan uint32)
	i2 := 100
	start := time.Now()

	go func() {
		for {
			select {
			case id := <-expire:
				delete(links, id)
				i2 = 100
			}
		}
	}()
	for true {
		send := pack(timer)

		i := 0
		for i == 0 {
			_, err := write(send, connection)
			if err != nil {
				log.Error(err)
			}

			// log.WithField("bytes", bytes).Info("wrote")

			time.Sleep(time.Duration(i2) / time.Duration(len(links)+1) * time.Millisecond)

			var id uint32
			_, i, id = read(connection, packRaceCounter(start))
			if i == 0 {
				i2 = i2 + 10
				log.WithField("interval", i2).Error("timeout")
			} else {
				if _, ok := links[id]; !ok {
					links[id] = link{
						id:     id,
						expire: expire,
						renew:  make(chan bool),
					}
					i2 = 100
					log.WithField("id", id).Info("New Car ")
					l := links[id]
					go l.timeout()
				}
				links[id].renew <- true
			}
		}
		// time.Sleep(100 * time.Millisecond)

	}

	time.Sleep(1 * time.Hour)

}

type link struct {
	id     uint32
	expire chan uint32
	renew  chan bool
}

func (l *link) timeout() {
	for {
		select {
		case <-time.After(2 * time.Second):
			l.expire <- l.id
			log.WithField("id", l.id).Info("Car timeout")
			return
		case <-l.renew:
			log.WithField("id", l.id).Info("Renew")
		}
	}
}

func packRaceCounter(start time.Time) []byte {
	centiSeconds := time.Now().Sub(start).Milliseconds() / 10
	be := make([]byte, 8)
	binary.BigEndian.PutUint64(be, uint64(centiSeconds))
	return be[len(be)-3:]
}

func read(connection io.ReadWriteCloser, timer []byte) (error, int, uint32) {
	buffer := make([]byte, 52)
	r := io.LimitReader(connection, 52)
	n := 0
	var err error

	n, err = r.Read(buffer)
	if err != nil {
		log.Error(err)
		return err, n, 0
	}
	/*	log.WithField("message", fmt.Sprintf("%x", buffer)).
		WithField("bytes", n).
		WithField("time", time.Now()).
		Info("received message from dongle")*/

	timer = buffer[9:12]
	return nil, n, uint32(buffer[1])
}

func write(send []byte, connection io.ReadWriteCloser) (int, error) {
	/*	log.WithFields(map[string]interface{}{
		"message": fmt.Sprintf("%x", send),
		"size":    fmt.Sprintf("%d", len(send)),
	}).Trace("send message to dongle")*/

	len, err := connection.Write(send)
	if err != nil {
		panic(err)
	}
	return len, nil
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
