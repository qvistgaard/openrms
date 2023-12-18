package webserver

import "sync"

type Event struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

type WebServer interface {
	RunServer(*sync.WaitGroup)
	PublishEvent(Event)
}
