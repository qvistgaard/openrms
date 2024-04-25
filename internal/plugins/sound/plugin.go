package sound

import (
	"embed"
	"github.com/gopxl/beep"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/plugins/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/ontrack"
	"github.com/qvistgaard/openrms/internal/plugins/pit"
	race2 "github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/sounds"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/system"
	"github.com/qvistgaard/openrms/internal/plugins/telemetry"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/rs/zerolog"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:embed announcements/finished.txt
//go:embed announcements/leading.txt
//go:embed announcements/limbmode.txt
//go:embed announcements/offtrack.txt
//go:embed announcements/out_of_fuel.txt
//go:embed announcements/pit_stop_complete.txt
//go:embed announcements/ready.txt
//go:embed announcements/result.txt
//go:embed announcements/start.txt
//go:embed announcements/fastest_lap.txt
//go:embed announcements/race_fastest_lap.txt
var announcements embed.FS

type Plugin struct {
	config       *system.Config
	sound        *system.Sound
	confirmation *confirmation.Plugin
	limbmode     *limbmode.Plugin
	fuel         *fuel.Plugin
	pit          *pit.Plugin
	tracker      tracker
	ontrack      *ontrack.Plugin
	telemetry    *telemetry.Plugin
	logger       zerolog.Logger
}

type tracker struct {
	raceState         race.Status
	ontrackCancel     map[types.CarId]chan bool
	cars              map[types.CarId]*car.Car
	maxDuration       time.Duration
	finalRise         *streamer.Playback
	finalRiseDuration time.Duration
	finalRisePlaying  bool
	duration          time.Duration
}

func New(logger zerolog.Logger, config *system.Config, sound *system.Sound, telemetry *telemetry.Plugin, race *race.Race, confirmation *confirmation.Plugin, limbMode *limbmode.Plugin, fuel *fuel.Plugin, pit *pit.Plugin, ontrack *ontrack.Plugin, plugin *race2.Plugin) (*Plugin, error) {
	p := &Plugin{
		logger:       logger.Level(config.Plugin.Sound.LogLevel),
		config:       config,
		sound:        sound,
		confirmation: confirmation,
		limbmode:     limbMode,
		fuel:         fuel,
		pit:          pit,
		ontrack:      ontrack,
		telemetry:    telemetry,
		tracker: tracker{
			ontrackCancel: make(map[types.CarId]chan bool),
			cars:          make(map[types.CarId]*car.Car),
		},
	}

	if p.config.Plugin.Sound.Enabled {
		p.registerObservers(race, confirmation, plugin)
	}
	return p, nil
}

func (p *Plugin) registerObservers(r *race.Race, confirmation *confirmation.Plugin, racePlugin *race2.Plugin) {
	r.Status().RegisterObserver(func(status race.Status) {
		if status == race.Stopped && p.tracker.raceState == race.Running {
			p.postRaceSequence()
		}

		if status == race.Running && p.tracker.raceState == race.Stopped && p.config.Plugin.Sound.Effects.Announcements.AfterStart {
			p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/start.txt",
				Random:   true,
			})
		}
		p.tracker.raceState = status
	})

	confirmation.Active().RegisterObserver(func(b bool) {
		if b && p.config.Plugin.Sound.Effects.Announcements.BeforeStart {
			p.sound.StopMusic()
			p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/ready.txt",
				Random:   true,
			})
		}
	})

	racePlugin.MaxDuration().RegisterObserver(func(duration time.Duration) {
		p.tracker.maxDuration = duration
	})

	r.Duration().RegisterObserver(func(duration time.Duration) {
		if p.tracker.raceState == race.Running && p.tracker.maxDuration > 0 {
			p.tracker.duration = duration

			if p.config.Plugin.Sound.Effects.Music.PreRaceFinish {
				timeUntilCompletion := p.tracker.maxDuration - duration
				if p.tracker.finalRise == nil {
					p.tracker.finalRise = sounds.EpicRise()
					p.tracker.finalRiseDuration = p.tracker.finalRise.SoftLenAsDuration()
					p.tracker.finalRisePlaying = false
				}

				if timeUntilCompletion <= p.tracker.finalRiseDuration && !p.tracker.finalRisePlaying {
					p.tracker.finalRisePlaying = true
					p.tracker.finalRise.SeekToPositionInDuration(p.tracker.finalRiseDuration - timeUntilCompletion)
					p.tracker.finalRise.Mute()
					p.sound.PlayMusic(p.tracker.finalRise, func() {
						p.tracker.finalRise = nil
					})
					p.tracker.finalRise.FadeIn(5 * time.Second)
				}
			}
		}
	})

	p.telemetry.Leader().RegisterObserver(func(id types.CarId) {
		if p.config.Plugin.Sound.Effects.Announcements.NewLeader {
			if p.tracker.raceState == race.Running && p.tracker.duration > 1*time.Minute {
				p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
					Fs:       announcements,
					Filename: "announcements/leading.txt",
					Random:   true,
					Data:     p.tracker.cars[id].TemplateData(),
				})
			} else {
				p.logger.Info().Msg("Leader updated within the first minute of the race. ignoring")
			}
		}
	})

	p.telemetry.FastestLap().RegisterObserver(func(id types.CarId) {
		if p.config.Plugin.Sound.Effects.Announcements.FastestLap {
			if p.tracker.raceState == race.Running && p.tracker.duration > 1*time.Minute {
				p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
					Fs:       announcements,
					Filename: "announcements/race_fastest_lap.txt",
					Random:   true,
					Data:     p.tracker.cars[id].TemplateData(),
				})
			} else {
				p.logger.Info().Msg("fastest lap updated within the first minute of the race. ignoring")
			}
		}
	})

}

func (p *Plugin) ConfigureCar(car *car.Car) {
	if p.config.Plugin.Sound.Enabled {
		car.LastLap().RegisterObserver(func(lap types.Lap) {
			if p.tracker.raceState == race.Running {
				if lap.Number > 0 {
					p.logger.Info().Msg("Playing Lap sound")
					playback := sounds.Lap()
					p.sound.PlayEffect(beep.Seq(playback, beep.Callback(func() {
						playback.Close()
					})))
				}

				u, err := p.fuel.Fuel(car.Id())
				if err != nil {
					log.Error(err)
					return
				}
				a, err := p.fuel.Average(car.Id())
				if err != nil {
					log.Error(err)
					return
				}
				f := u.Get() / a

				if a > 0 && f < 5 && p.config.Plugin.Sound.Effects.Announcements.OutOfFuel {
					p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
						Fs:       announcements,
						Filename: "announcements/out_of_fuel.txt",
						Random:   true,
						Data:     car.TemplateData(),
					})
				}
			}
		})

		if p.config.Plugin.Sound.Announcements.Enabled {
			p.tracker.cars[car.Id()] = car
			p.limbmode.LimbMode(car.Id()).RegisterObserver(func(b bool) {
				if b && p.config.Plugin.Sound.Effects.Announcements.LimbMode {
					p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
						Fs:       announcements,
						Filename: "announcements/limbmode.txt",
						Random:   true,
						Data:     car.TemplateData(),
					})
				}
			})

			p.pit.Active(car.Id()).RegisterObserver(func(b bool) {
				if !b && p.config.Plugin.Sound.Effects.Announcements.PitStopComplete {
					p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
						Fs:       announcements,
						Filename: "announcements/pit_stop_complete.txt",
						Random:   true,
						Data:     car.TemplateData(),
					})
				}
			})

			p.ontrack.Ontrack(car.Id()).RegisterObserver(func(b bool) {
				if !b && p.tracker.raceState == race.Running && p.config.Plugin.Sound.Effects.Announcements.OffTrack {
					p.sound.OptionalAnnouncement(&announcer.ReadFileTemplateAnnouncement{
						Fs:       announcements,
						Filename: "announcements/offtrack.txt",
						Random:   true,
						Data:     car.TemplateData(),
					})
				}
			})
		}
	}

}

func (p *Plugin) InitializeCar(car *car.Car) {

}

func (p *Plugin) Priority() int {
	return 1000
}

func (p *Plugin) Name() string {
	return "sound"
}
