package messages

import (
	"github.com/qvistgaard/openrms/internal/plugins/telemetry"
	"github.com/qvistgaard/openrms/internal/state/race"
	"time"
)

type RaceDuration time.Duration

type Update struct {
	RaceTelemetry telemetry.Race
	RaceStatus    race.Status
	RaceDuration  time.Duration
	TrackMaxSpeed uint8
}
