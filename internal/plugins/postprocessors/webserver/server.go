package webserver

import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	/*	race       chan state.CourseState
		car        chan state.CarState*/
	listen  string
	command chan<- interface{}
	context *application.Context
}

type StateMessage struct {
	/*	Cars   []state.CarState    `json:"cars"`
		Course []state.CourseState `json:"race"`*/
}

/*func (s *Server) CarChannel() chan<- state.CarState {
	return s.car
}

func (s *Server) RaceChannel() chan<- state.CourseState {
	return s.race
}*/

func (s *Server) CommandChannel(c chan<- interface{}) {
	s.command = c
}

func (s *Server) Process() {
	defer func() {
		log.Fatal("Websocket process died")
	}()

	go s.processWebsocket()

	http.HandleFunc("/static/", web.StaticContentHandler)
	http.Handle("/", http.RedirectHandler("/static/index.html", 301))

	log.Infof("started webserver post processor, listening on %s", s.listen)
	err := http.ListenAndServe(s.listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Warn("webserver post processor stopped")
}

func (s *Server) processWebsocket() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			/*		case message := <-s.:
					for client := range s.clients {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(s.clients, client)
						}
					}*/
			/*		case message := <-s.race:
					for client := range s.clients {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(s.clients, client)
						}
					}*/
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
