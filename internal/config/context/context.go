package context

import (
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state"
)

type Config map[string]interface{}
type Context struct {
	Config         *Config
	Implement      implement.Implementer
	Cars           car.Repository
	Rules          state.Rules
	Postprocessors postprocess.PostProcess
	Course         *state.Course
}
