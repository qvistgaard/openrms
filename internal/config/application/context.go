package application

import (
	"github.com/madflojo/tasks"
	"github.com/qvistgaard/openrms/internal/implement"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/leaderboard"
	"github.com/qvistgaard/openrms/internal/postprocess"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/rules"
	"github.com/qvistgaard/openrms/internal/state/track"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/qvistgaard/openrms/internal/webserver"
)

type Config map[string]interface{}
type Context struct {
	Config         *Config
	Implement      implement.Implementer
	Cars           car.Repository
	Rules          rules.Rules
	Postprocessors *postprocess.PostProcess
	Scheduler      *tasks.Scheduler
	Webserver      webserver.WebServer
	Track          *track.Track
	Race           *race.Race
	ValueFactory   *reactive.Factory
	Leaderboard    *leaderboard.Leaderboard
}
