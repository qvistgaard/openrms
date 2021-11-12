package webserver

type Event struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

type WebServer interface {
	RunServer()
	PublishEvent(Event)
}
