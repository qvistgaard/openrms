package websocket

import (
	"encoding/json"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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
	Cars []state.CarChanges    `json:"cars"`
	Race []state.CourseChanges `json:"race"`
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
	var stateMessages = StateMessage{
		Cars: []state.CarChanges{},
		Race: []state.CourseChanges{},
	}

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
			stateMessages.Cars = append(stateMessages.Cars, message)

		case <-time.After(500 * time.Millisecond):
			if len(stateMessages.Cars) > 0 || len(stateMessages.Race) > 0 {
				marshal, _ := json.Marshal(stateMessages)
				for client := range b.clients {
					select {
					case client.send <- marshal:
					default:
						close(client.send)
						delete(b.clients, client)
					}
				}
			}
			stateMessages = StateMessage{
				Cars: []state.CarChanges{},
				Race: []state.CourseChanges{},
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

	client := &Client{broadcast: ws, conn: conn, send: make(chan []byte, 256), request: r}
	client.broadcast.register <- client
	log.WithFields(map[string]interface{}{
		"ip_addr": conn.RemoteAddr(),
	}).Infof("Client connected")

	go client.writePump()
	go client.read()
}
