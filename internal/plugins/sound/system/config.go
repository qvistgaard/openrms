package system

import (
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/elevenlabs"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/playht"
	"github.com/rs/zerolog"
)

type Config struct {
	Plugin struct {
		Sound struct {
			LogLevel zerolog.Level `mapstructure:"loglevel"`
			Enabled  bool          `mapstructure:"enabled" default:"true"`
			Effects  struct {
				Announcements *Announcements
				Music         *Music
			}
			Announcements struct {
				Enabled    bool                         `mapstructure:"enabled" default:"true"`
				Engine     string                       `mapstructure:"engine" default:"playht"`
				PlayHT     *playht.PlayHTConfig         `mapstructure:"playht"`
				ElevenLabs *elevenlabs.ElevenLabsConfig `mapstructure:"elevenlabs"`
			} `mapstructure:"announcements"`
		} `mapstructure:"sound"`
	} `mapstructure:"plugins"`
}

type Announcements struct {
	BeforeStart     bool `mapstructure:"before-start" default:"false"`
	AfterStart      bool `mapstructure:"after-start" default:"false"`
	NewLeader       bool `mapstructure:"new-leader" default:"true"`
	FastestLap      bool `mapstructure:"fastest-lap" default:"true"`
	OutOfFuel       bool `mapstructure:"out-of-fuel" default:"true"`
	LimbMode        bool `mapstructure:"limb-mode" default:"true"`
	PitStopComplete bool `mapstructure:"pit-stop-complete" default:"true"`
	OffTrack        bool `mapstructure:"off-track" default:"false"`
}

type Music struct {
	PreRaceFinish bool `mapstructure:"pre-race-finish" default:"false"`
	PostRace      bool `mapstructure:"post-race" default:"true"`
}
