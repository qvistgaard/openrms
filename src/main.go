package main

import (
	"./plugins/connector/oxigen"
	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/tarm/serial"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup // 1

func main() {
	var err error
	// connection, _ := oxigen.NewUSBConnection("COM5")

	c := &serial.Config{Name: "COM5", Baud: 1200, ReadTimeout: time.Millisecond * 1000}
	connection, _ := serial.OpenPort(c)
	o, err := oxigen.Connect(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Closer()

	input := queue.NewFIFO()
	output := queue.NewFIFO()

	// o.Start()

	wg.Add(1)
	go eventloop(o, input, output)
	// go o.EventLoop(input, output)
	/*
		for {
			elm, err := output.DequeueOrWaitForNextElement()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%x", elm)
		}
	*/
	wg.Wait()
}

func eventloop(o *oxigen.Oxigen, input queue.Queue, output queue.Queue) error {
	defer wg.Done()
	return o.EventLoop(input, output)
}
