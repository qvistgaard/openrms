package main

import (
	"./plugins/connector/oxigen"
	queue "github.com/enriquebris/goconcurrentqueue"
	"log"
	"sync"
)

var wg sync.WaitGroup // 1

func main() {
	var err error
	connection, _ := oxigen.NewUSBConnection("COM5")
	o, err := oxigen.Connect(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Closer()

	input := queue.NewFIFO()
	output := queue.NewFIFO()

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
