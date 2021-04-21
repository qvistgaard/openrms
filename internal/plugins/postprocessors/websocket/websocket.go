package websocket

import (
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WebSocket struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	race       chan state.CourseChanges
	car        chan state.CarChanges
	listen     string
}

type StateMessage struct {
	Cars   []state.CarChanges    `json:"cars"`
	Course []state.CourseChanges `json:"race"`
}

func (b *WebSocket) CarChannel() chan<- state.CarChanges {
	return b.car
}

func (b *WebSocket) RaceChannel() chan<- state.CourseChanges {
	return b.race
}

func (b *WebSocket) Process() {
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
		}
	}
}

func serveWebSocket(ws *WebSocket, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{broadcast: ws, conn: conn, send: make(chan interface{}, 256), request: r}
	client.broadcast.register <- client
	log.WithFields(map[string]interface{}{
		"ip_addr": conn.RemoteAddr(),
	}).Infof("Client connected")

	go client.writePump()
	go client.read()
}
