package limbmode

import "github.com/qvistgaard/openrms/internal/config/context"

func CreateFromConfig(ctx *context.Context) *LimbMode {
	return &LimbMode{
		MaxSpeed: 100,
	}
}
