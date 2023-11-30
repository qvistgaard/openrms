package messages

import (
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

type RaceDuration time.Duration

type Update struct {
	RaceTelemetry types.RaceTelemetry
	RaceStatus    race.RaceStatus
	RaceDuration  time.Duration
	TrackMaxSpeed uint8
}
