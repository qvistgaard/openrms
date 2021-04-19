package websocket

import (
	"github.com/qvistgaard/openrms/internal/state"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Postprocessors struct {
		WebSocket struct {
			Listen string
		}
	}
}

func CreateFromConfig(config []byte) (*WebSocket, error) {
	c := &Config{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}

	ws := &WebSocket{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		race:       make(chan state.RaceChanges, 1024),
		car:        make(chan state.CarChanges, 1024),
		listen:     c.Postprocessors.WebSocket.Listen,
	}

	return ws, nil
}
