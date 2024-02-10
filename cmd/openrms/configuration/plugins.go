package configuration

import (
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins"
	"github.com/qvistgaard/openrms/internal/plugins/commentary"
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/elevenlabs"
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/playht"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/ontrack"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	"github.com/qvistgaard/openrms/internal/plugins/race"
	race2 "github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/track"
)

func Plugins(conf Config) (*plugins.Plugins, error) {
	c := &plugins.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return plugins.New(c)
}

// FuelPlugin initializes and returns a new fuel plugin instance based on the provided configuration and LimpMode plugin.
// It takes a `Config` map containing the fuel plugin configuration settings and an optional `limbmode.Plugin` instance.
//
// The `conf` parameter should be a `Config` map containing the fuel plugin configuration settings.
//
// The `limpMode` parameter is an optional instance of the `limbmode.Plugin` type that can be provided if needed.
// If not needed, you can pass `nil` for this parameter.
//
// Returns:
//   - A new instance of the 'fuel.Plugin' type representing the initialized fuel plugin.
//   - An error if there was an issue initializing the fuel plugin instance.
func FuelPlugin(conf Config, limpMode *limbmode.Plugin, commentary *commentary.Plugin) (*fuel.Plugin, error) {
	c := &fuel.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}

	return fuel.New(*c, limpMode, commentary)
}

// RacePlugin initializes and returns a new race plugin instance based on the provided configuration.
// It takes a `Config` map containing the race plugin configuration settings.
//
// The `conf` parameter should be a `Config` map containing the race plugin configuration settings.
//
// Returns:
//   - A new instance of the 'race.Plugin' type representing the initialized race plugin.
//   - An error if there was an issue initializing the race plugin instance.
func RacePlugin(_ Config, r *race2.Race, plugin *confirmation.Plugin, comment *commentary.Plugin) (*race.Plugin, error) {
	return race.New(r, plugin, comment)
}

func CommentaryPlugin(conf Config) (*commentary.Plugin, error) {
	c := &commentary.Config{}
	defaults.SetDefaults(c)
	c.Plugin.Commentary.PlayHT = &playht.PlayHTConfig{}
	c.Plugin.Commentary.ElevenLabs = &elevenlabs.ElevenLabsConfig{}
	defaults.SetDefaults(c.Plugin.Commentary.PlayHT)
	defaults.SetDefaults(c.Plugin.Commentary.ElevenLabs)
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return commentary.New(c)
}

func ConfirmationPlugin(conf Config, commentary *commentary.Plugin) (*confirmation.Plugin, error) {
	c := &confirmation.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return confirmation.New(c, commentary)
}

func PitPlugin(conf Config, comment *commentary.Plugin, stops ...pit.SequencePlugin) (*pit.Plugin, error) {
	c := &pit.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return pit.New(c, comment, stops...)
}

func FlagPlugin(conf Config, track *track.Track, race *race2.Race) (*flags.Plugin, error) {
	c := &flags.Config{}
	defaults.SetDefaults(c)
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read flag plugin configuration")
	}
	return flags.New(c, track, race)
}

func OnTrackPlugin(conf Config, f *flags.Plugin, comment *commentary.Plugin) (*ontrack.Plugin, error) {
	c := &ontrack.Config{}
	defaults.SetDefaults(c)
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return ontrack.New(c, f, comment)
}

// LimbModePlugin initializes and returns a new LimpMode plugin instance based on the provided configuration.
// It takes a `Config` map containing the LimpMode plugin configuration settings.
//
// The `conf` parameter should be a `Config` map containing the LimpMode plugin configuration settings.
//
// Returns:
//   - A new instance of the 'limbmode.Plugin' type representing the initialized LimpMode plugin.
//   - An error if there was an issue initializing the LimpMode plugin instance.
func LimbModePlugin(conf Config, commentary *commentary.Plugin) (*limbmode.Plugin, error) {
	c := &limbmode.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return limbmode.New(c, commentary)
}
