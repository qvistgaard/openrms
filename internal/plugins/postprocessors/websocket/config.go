package websocket

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

type Config struct {
	Postprocessors struct {
		WebSocket struct {
			Listen string
		}
	}
}

func CreateFromConfig(ctx *context.Context) (*WebSocket, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	ws := &WebSocket{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		race:       make(chan state.CourseState, 1024),
		car:        make(chan state.CarState, 1024),
		context:    ctx,
		listen:     c.Postprocessors.WebSocket.Listen,
	}
	return ws, nil
}
