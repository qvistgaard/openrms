package main

import (
	"../plugins/connector/oxigen"
	queue "github.com/enriquebris/goconcurrentqueue"
	"log"
)

func main() {
	var err error
	usb := oxigen.NewUSBConnection("/dev/ttyACM1")
	o, err := oxigen.Connect(usb)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Closer()

	input := queue.NewFIFO()
	output := queue.NewFIFO()

	go o.EventLoop(input, output)

	for {
		elm, err := output.DequeueOrWaitForNextElement()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%x", elm)
	}
}
