package websocket

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WebSocket struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	race       chan state.CourseState
	car        chan state.CarState
	listen     string
	command    chan<- interface{}
	context    *context.Context
}

type StateMessage struct {
	Cars   []state.CarState    `json:"cars"`
	Course []state.CourseState `json:"race"`
}

func (b *WebSocket) CarChannel() chan<- state.CarState {
	return b.car
}

func (b *WebSocket) RaceChannel() chan<- state.CourseState {
	return b.race
}

func (b *WebSocket) CommandChannel(c chan<- interface{}) {
	b.command = c
}

func (b *WebSocket) Process() {
	defer func() {
		log.Fatal("Websocket process died")
	}()

	log.Infof("started websocket post processor, listening on %s", b.listen)

	go b.processWebsocket()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(b, w, r)
	})

	err := http.ListenAndServe(b.listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Warn("websocket post processor stopped")
}

func (b *WebSocket) processWebsocket() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
		case client := <-b.unregister:
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.send)
			}
		case message := <-b.car:
			for client := range b.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(b.clients, client)
				}
			}
		case message := <-b.race:
			for client := range b.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(b.clients, client)
				}
			}
		}
	}
}

func serveWebSocket(ws *WebSocket, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		broadcast: ws,
		conn:      conn,
		send:      make(chan interface{}, 256),
		request:   r,
		command:   ws.command,
		context:   ws.context,
	}
	client.broadcast.register <- client
	log.WithFields(map[string]interface{}{
		"ip_addr": conn.RemoteAddr(),
	}).Infof("Client connected")

	go client.writePump()
	go client.read()
}
