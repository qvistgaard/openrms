package race

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

type Rule interface {
	ConfigureRaceState(*race.Race)
	InitializeRaceState(*race.Race, context.Context, reactive.ValuePostProcessor)
}
