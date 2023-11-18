package config

import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

func CreateValueFactory(context *application.Context) error {
	context.ValueFactory = reactive.NewFactory(context.Postprocessors.ValuePostProcessor())
	return nil
}
