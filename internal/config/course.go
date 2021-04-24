package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
)

func CreateCourse(ctx *context.Context) error {
	c := &state.CourseConfig{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil
	}
	ctx.Course = state.CreateCourse(c, ctx.Rules)
	return nil
}
