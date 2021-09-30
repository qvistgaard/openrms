package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/influxdb"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/webserver"
	"github.com/qvistgaard/openrms/internal/postprocess"
)

type PostProcessorConfig struct {
	Postprocessors map[string]interface{}
}

func CreatePostProcessors(context *context.Context) error {
	c := &PostProcessorConfig{}
	err := mapstructure.Decode(context.Config, c)

	if err != nil {
		return err
	}

	var postprocessors []postprocess.PostProcessor
	for k, _ := range c.Postprocessors {
		switch k {
		case "influxdb":
			p, err := influxdb.CreateFromConfig(context)
			if err != nil {
				return err
			}
			postprocessors = append(postprocessors, p)
			go p.Process()
		case "webserver":
			ws, err := webserver.CreateFromConfig(context)
			if err != nil {
				return err
			}
			postprocessors = append(postprocessors, ws)
			go ws.Process()
		}
	}
	context.Postprocessors = postprocess.CreatePostProcess(postprocessors)
	return nil
}
