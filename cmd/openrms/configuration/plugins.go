package configuration

import (
	"context"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/ontrack"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	"github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/plugins/sound"
	ell "github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/elevenlabs"
	pht "github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/playht"
	"github.com/qvistgaard/openrms/internal/plugins/sound/system"
	"github.com/qvistgaard/openrms/internal/plugins/telemetry"
	race2 "github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/track"
)

// Plugins initializes and returns a new Plugins instance based on the provided configuration.
//
// Parameters:
// - conf: A Config map containing configuration settings.
//
// Returns:
//   - A pointer to a newly created Plugins instance.
//   - An error if there is an issue decoding the configuration or initializing the Plugins instance,
//     including detailed message about the failure.
//
// This function decodes the provided configuration into a Plugins Config structure and uses it
// to create a new Plugins instance. It's a general-purpose initializer for any Plugins type,
// handling configuration errors with descriptive messages.
func Plugins(conf Config) (*plugins.Plugins, error) {
	c := &plugins.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return plugins.New(c)
}

// FuelPlugin initializes and returns a new fuel plugin instance based on the provided configuration.
// It integrates with an optional LimpMode plugin and a SoundFX plugin for enhanced functionality.
//
// Parameters:
// - conf: A Config map containing the fuel plugin configuration settings.
// - limpMode: An optional *limbmode.Plugin instance for integrating limp mode functionality.
//             Pass nil if limp mode integration is not needed.
// - sound: A *sound.Plugin instance for sound effects integration.
//
// Returns:
// - A new *fuel.Plugin instance configured according to the provided settings.
// - An error if there is an issue decoding the configuration or initializing the fuel plugin instance,
//   with a detailed error message.
//
// This function decodes the configuration into a fuel.Config structure, then initializes and
// returns a fuel plugin instance with optional limp mode and sound effects functionality.

func FuelPlugin(conf Config, limpMode *limbmode.Plugin, sound *system.Sound) (*fuel.Plugin, error) {
	c := &fuel.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}

	return fuel.New(*c, limpMode, sound)
}

// RacePlugin initializes and returns a new race plugin instance based on the provided configuration.
// It configures the plugin with a race instance, a confirmation plugin, and a sound effects plugin.
//
// Parameters:
// - conf: A Config map containing the race plugin configuration settings. (Currently unused placeholder)
// - r: A *race2.Race instance to be associated with the race plugin.
// - plugin: A *confirmation.Plugin instance for race confirmation functionalities.
// - sound: A *sound.Plugin instance for integrating sound effects into the race plugin.
//
// Returns:
// - A new *race.Plugin instance configured with the provided race, confirmation, and sound plugins.
// - An error if there is an issue initializing the race plugin instance.
//
// This function creates a race plugin instance, leveraging the given race mechanics, confirmation logic,
// and audio effects for a comprehensive race management solution.
func RacePlugin(_ Config, r *race2.Race, plugin *confirmation.Plugin, sound *system.Sound) (*race.Plugin, error) {
	return race.New(r, plugin, sound)
}

// ConfirmationPlugin initializes and returns a new confirmation plugin instance based on the provided
// configuration. It integrates with a SoundFX plugin for audio feedback on confirmations.
//
// Parameters:
// - conf: A Config map containing the confirmation plugin configuration settings.
// - sound: A *sound.Plugin instance for integrating sound effects into confirmation feedback.
//
// Returns:
//   - A new *confirmation.Plugin instance configured according to the provided settings.
//   - An error if there is an issue decoding the configuration or initializing the confirmation plugin instance,
//     with a detailed error message.
//
// This function decodes the configuration into a confirmation.Config structure, then initializes and
// returns a confirmation plugin instance with sound effects functionality for audio feedback.
func ConfirmationPlugin(conf Config) (*confirmation.Plugin, error) {
	c := &confirmation.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return confirmation.New(c)
}

// PitPlugin initializes and returns a new pit plugin instance based on the provided configuration.
// It integrates with a SoundFX plugin and optional sequence plugins for detailed pit stop management.
//
// Parameters:
// - conf: A Config map containing the pit plugin configuration settings.
// - sound: A *sound.Plugin instance for integrating sound effects into pit stop events.
// - stops: Variadic parameters of pit.SequencePlugin instances for defining pit stop sequences.
//
// Returns:
//   - A new *pit.Plugin instance configured with sound effects and optional pit stop sequences.
//   - An error if there is an issue decoding the configuration or initializing the pit plugin instance,
//     with a detailed error message.
//
// This function decodes the configuration into a pit.Config structure, then initializes and
// returns a pit plugin instance, providing sound effects for pit events and customizable pit stop sequences.
func PitPlugin(conf Config, sound *system.Sound, stops ...pit.SequencePlugin) (*pit.Plugin, error) {
	c := &pit.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return pit.New(c, sound, stops...)
}

// FlagPlugin initializes and returns a new flags plugin instance based on the provided configuration.
// It integrates with track and race instances for comprehensive flag management during races.
//
// Parameters:
// - conf: A Config map containing the flags plugin configuration settings.
// - track: A *track.Track instance for track-related flag functionalities.
// - race: A *race2.Race instance for race-specific flag management.
//
// Returns:
//   - A new *flags.Plugin instance configured according to the provided settings.
//   - An error if there is an issue decoding the configuration or initializing the flags plugin instance,
//     with a detailed error message.
//
// This function decodes the configuration into a flags.Config structure, sets default values,
// then initializes and returns a flags plugin instance, integrating track and race contexts for flag management.
func FlagPlugin(conf Config, track *track.Track, race *race2.Race) (*flags.Plugin, error) {
	c := &flags.Config{}
	defaults.SetDefaults(c)
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read flag plugin configuration")
	}
	return flags.New(c, track, race)
}

// OnTrackPlugin initializes and returns a new on-track plugin instance based on the provided configuration.
// It integrates with a Flags plugin and a SoundFX plugin for enhanced track management and audio feedback.
//
// Parameters:
// - conf: A Config map containing the on-track plugin configuration settings.
// - f: A *flags.Plugin instance for integrating flag management functionalities.
// - sound: A *sound.Plugin instance for sound effects integration.
//
// Returns:
//   - A new *ontrack.Plugin instance configured according to the provided settings.
//   - An error if there is an issue decoding the configuration or initializing the on-track plugin instance,
//     with a detailed error message.
//
// This function decodes the configuration into an ontrack.Config structure, sets default values,
// then initializes and returns an on-track plugin instance, providing enhanced track management and audio feedback.
func OnTrackPlugin(conf Config, f *flags.Plugin, sound *system.Sound) (*ontrack.Plugin, error) {
	c := &ontrack.Config{}
	defaults.SetDefaults(c)
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return ontrack.New(c, f, sound)
}

// SoundPlugin initializes and returns a new SoundFX plugin instance based on the provided configuration.
// It configures the plugin with context-specific sound effects and announcement configurations.
//
// Parameters:
// - ctx: A context.Context instance for managing the plugin's lifecycle.
// - conf: A Config map containing the SoundFX plugin configuration settings.
//
// Returns:
//   - A new *sound.Plugin instance configured with advanced sound effects and announcement capabilities.
//   - An error if there is an issue decoding the configuration or initializing the SoundFX plugin instance,
//     with a detailed error message.
//
// This function sets default values for the SoundFX configuration, adjusts specific announcement settings,
// and returns a SoundFX plugin instance tailored to the provided configuration, enhancing the audio experience.
func SoundPlugin(conf Config, soundSystem *system.Sound, telemetry *telemetry.Plugin, race *race2.Race, confirmation *confirmation.Plugin, limbMode *limbmode.Plugin, fuel *fuel.Plugin, pit *pit.Plugin, ontrack *ontrack.Plugin, plugin *race.Plugin) (*sound.Plugin, error) {
	c := &system.Config{}
	defaults.SetDefaults(c)

	c.Plugin.Sound.Announcements.PlayHT = &pht.PlayHTConfig{}
	c.Plugin.Sound.Announcements.ElevenLabs = &ell.ElevenLabsConfig{}
	defaults.SetDefaults(c.Plugin.Sound.Announcements.PlayHT)
	defaults.SetDefaults(c.Plugin.Sound.Announcements.ElevenLabs)

	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read sound plugin configuration")
	}

	return sound.New(c, soundSystem, telemetry, race, confirmation, limbMode, fuel, pit, ontrack, plugin)
}

// LimbModePlugin initializes and returns a new LimpMode plugin instance based on the provided configuration.
// It integrates with a SoundFX plugin for audio feedback on limp mode activations.
//
// Parameters:
// - conf: A Config map containing the LimpMode plugin configuration settings.
// - sound: A *sound.Plugin instance for integrating sound effects into limp mode feedback.
//
// Returns:
// - A new *limbmode.Plugin instance configured according to the provided settings.
// - An error if there is an issue decoding the configuration or initializing the LimpMode plugin instance,
//   with a detailed error message.
//
// This function decodes the configuration into a limbmode.Config structure, then initializes and
// returns a LimpMode plugin instance, providing sound effects for limp mode activations and feedback.

func LimbModePlugin(conf Config, sound *system.Sound) (*limbmode.Plugin, error) {
	c := &limbmode.Config{}
	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read fuel plugin configuration")
	}
	return limbmode.New(c, sound)
}

func SoundSystem(ctx context.Context, conf Config) (*system.Sound, error) {
	c := &system.Config{}
	defaults.SetDefaults(c)

	c.Plugin.Sound.Announcements.PlayHT = &pht.PlayHTConfig{}
	c.Plugin.Sound.Announcements.ElevenLabs = &ell.ElevenLabsConfig{}
	defaults.SetDefaults(c.Plugin.Sound.Announcements.PlayHT)
	defaults.SetDefaults(c.Plugin.Sound.Announcements.ElevenLabs)

	err := mapstructure.Decode(conf, c)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read sound plugin configuration")
	}

	return system.New(ctx, c)
}
