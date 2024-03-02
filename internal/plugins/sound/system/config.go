package system

import (
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/elevenlabs"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/playht"
)

type Config struct {
	Plugin struct {
		Sound struct {
			Enabled       bool `mapstructure:"enabled" default:"true"`
			Announcements struct {
				Enabled    bool                         `mapstructure:"enabled" default:"true"`
				Engine     string                       `mapstructure:"engine" default:"playht"`
				PlayHT     *playht.PlayHTConfig         `mapstructure:"playht"`
				ElevenLabs *elevenlabs.ElevenLabsConfig `mapstructure:"elevenlabs"`
			} `mapstructure:"announcements"`
		} `mapstructure:"sound"`
	} `mapstructure:"plugins"`
}
