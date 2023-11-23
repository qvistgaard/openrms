package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/race"
)

type Rule interface {
	ConfigureRaceState(*race.Race)
	InitializeRaceState(*race.Race, context.Context)
}
