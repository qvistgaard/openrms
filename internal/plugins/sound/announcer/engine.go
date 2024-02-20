package announcer

import (
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
)

type Engine interface {
	Announce(announcement Announcement) (*streamer.Playback, error)
}
