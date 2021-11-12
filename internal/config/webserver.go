package config

import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/plugins/webserver"
)

func CreateWebserver(ctx *application.Context) error {
	var err error
	ctx.Webserver, err = webserver.CreateFromConfig(ctx)
	return err
}
