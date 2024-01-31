package commentary

import (
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/elevenlabs"
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/playht"
)

type Config struct {
	Plugin struct {
		Commentary struct {
			Enabled    bool                         `mapstructure:"enabled" default:"true"`
			Engine     string                       `mapstructure:"engine" default:"playht"`
			PlayHT     *playht.PlayHTConfig         `mapstructure:"playht"`
			ElevenLabs *elevenlabs.ElevenLabsConfig `mapstructure:"elevenlabs"`
		} `mapstructure:"commentary"`
	} `mapstructure:"plugins"`
}
