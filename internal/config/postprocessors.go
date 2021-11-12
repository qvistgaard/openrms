package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/influxdb"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/leaderboard"
	"github.com/qvistgaard/openrms/internal/postprocess"
)

type PostProcessorConfig struct {
	Postprocessors map[string]interface{}
}

func CreatePostProcessors(context *application.Context) error {
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
		case "leaderboard":
			lb, err := leaderboard.CreateFromConfig(context)
			if err != nil {
				return err
			}
			postprocessors = append(postprocessors, lb)
		}
	}
	context.Postprocessors = postprocess.CreatePostProcess(postprocessors)
	return nil
}
