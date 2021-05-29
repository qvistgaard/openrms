package websocket

import (
	"encoding/json"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Command struct {
	Race *struct {
		Name  string
		Value interface{}
	}
	Car *struct {
		CarId int64
		Name  string
		Value interface{}
	}
	Get *struct {
		Car *struct {
			CarId []state.CarId
			Name  []string
		}
		Race *struct {
			Name []string
		}
	}
}

type State struct {
	Race map[string]interface{}                 `json:"race"`
	Cars map[state.CarId]map[string]interface{} `json:"cars"`
}

// Client is a middleman between the websocket connection and the broadcast.
type Client struct {
	broadcast *WebSocket

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send    chan interface{}
	request *http.Request
	command chan<- interface{}
	context *context.Context
}

// read pumps messages from the websocket connection to the broadcast.
//
// The application runs read in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) read() {
	defer func() {
		c.broadcast.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		o := &Command{}
		ok := json.Unmarshal(b, o)
		if ok == nil {
			if o.Race != nil {
				log.WithField("state", o.Race.Name).
					WithField("value", o.Race.Value).
					Infof("race command recieved")
				c.command <- state.CourseCommand{Name: o.Race.Name, Value: o.Race.Value}
			}
			if o.Car != nil {
				log.WithField("state", o.Car.Name).
					WithField("value", o.Car.Value).
					WithField("car", o.Car.CarId).
					Infof("car command received")
					//if cid, ok := o.Car.CarId.Int64(); ok == nil {
				c.command <- state.CarCommand{
					CarId: state.CarId(o.Car.CarId),
					Name:  o.Car.Name,
					Value: o.Car.Value,
				}
				// }
			}
			if o.Get != nil {
				s := State{}
				if o.Get.Car != nil {
					s.Cars = make(map[state.CarId]map[string]interface{})
					for _, car := range o.Get.Car.CarId {
						s.Cars[car] = make(map[string]interface{})
						if c.context.Cars.Exists(car) {
							cs, _, _ := c.context.Cars.Get(car)
							for _, n := range o.Get.Car.Name {
								s.Cars[car][n] = cs.Get(n)
							}
						}
					}
				}
				if o.Get.Race != nil {
					s.Race = make(map[string]interface{})
					for _, n := range o.Get.Race.Name {
						s.Race[n] = c.context.Course.Get(n)
					}
				}
				marshal, _ := json.Marshal(s)
				c.send <- marshal
			}
		} else {
			log.Warn(ok)
			marshal, _ := json.Marshal(map[string]interface{}{
				"error": "no such command",
			})
			c.send <- marshal
		}
	}
}

// writePump pumps messages from the broadcast to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	var stateMessages = StateMessage{
		Cars:   []state.CarState{},
		Course: []state.CourseState{},
	}
	nextTX := time.Now()

	for _, car := range c.context.Cars.All() {
		carState := car.State()
		if c.filterCarChanges(carState) {
			stateMessages.Cars = append(stateMessages.Cars, carState)
		}
	}
	stateMessages.Course = append(stateMessages.Course, c.context.Course.State())

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The broadcast closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if carChanges, ok := message.(state.CarState); ok && c.filterCarChanges(carChanges) {
				stateMessages.Cars = append(stateMessages.Cars, carChanges)
			}
			if courseChanges, ok := message.(state.CourseState); ok {
				stateMessages.Course = append(stateMessages.Course, courseChanges)
			}
			if json, ok := message.([]byte); ok {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(json)
			}

			if time.Now().After(nextTX) {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if len(stateMessages.Cars) > 0 || len(stateMessages.Course) > 0 {
					w, err := c.conn.NextWriter(websocket.TextMessage)
					if err != nil {
						return
					}
					marshal, _ := json.Marshal(stateMessages)
					w.Write(marshal)
					stateMessages = StateMessage{
						Cars:   []state.CarState{},
						Course: []state.CourseState{},
					}
					if err := w.Close(); err != nil {
						log.Error("Conncetion error: " + err.Error())
						return
					}
				}
				nextTX = time.Now().Add(500 * time.Millisecond)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) collectAndSendChanges() {

}

func (c *Client) filterCarChanges(changes state.CarState) bool {
	get := c.request.URL.Query().Get("car")
	if get == "" || get == strconv.FormatUint(uint64(changes.Car), 10) {
		return true
	}
	return false
}
