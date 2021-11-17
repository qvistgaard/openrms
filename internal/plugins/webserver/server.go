package webserver

import (
	"encoding/json"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/webserver"
	"github.com/qvistgaard/openrms/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	listen     string
	events     chan webserver.Event
	context    *application.Context
}

func NewServer(listen string, context *application.Context) *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		listen:     listen,
		events:     make(chan webserver.Event, 100),
		context:    context,
	}
}

func (s *Server) PublishEvent(event webserver.Event) {
	s.events <- event
}

func (s *Server) RunServer() {
	defer func() {
		log.Fatal("Webserver died")
	}()

	go s.processWebsocket()

	http.HandleFunc("/static/", web.StaticContentHandler)
	http.Handle("/", http.RedirectHandler("/static/index.html", 301))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(s, w, r)
	})

	log.Infof("started webserver post processor, listening on %s", s.listen)
	err := http.ListenAndServe(s.listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Warn("webserver post processor stopped")
}

func (s *Server) processWebsocket() {
	defer func() {
		log.Fatal("Websocket processor died")
	}()
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}

		// Read from broadcast channel and send to client
		case message := <-s.events:
			marshal, _ := json.Marshal(message)
			for client := range s.clients {
				select {
				case client.send <- marshal:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

func serveWebSocket(ws *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		broadcast: ws,
		conn:      conn,
		send:      make(chan []byte, 256),
		request:   r,
		context:   ws.context,
	}
	client.broadcast.register <- client
	log.WithFields(map[string]interface{}{
		"ip_addr": conn.RemoteAddr(),
	}).Infof("Client connected")

	go client.writePump()
	go client.read()
}
