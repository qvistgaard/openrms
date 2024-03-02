package system

import (
	"context"
	"errors"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/elevenlabs"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer/engines/playht"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

type Sound struct {
	announcer   announcer.Engine
	config      *Config
	queue       chan beep.Streamer
	context     context.Context
	nowPlaying  *streamer.Playback
	effectCount atomic.Int64
}

func New(context context.Context, config *Config) (*Sound, error) {
	p := &Sound{
		config:  config,
		queue:   make(chan beep.Streamer, 20),
		context: context,
	}

	err := initCommentaryEngine(config, p)
	if err != nil {
		return nil, err
	}

	err = p.initSpeaker()
	if err != nil {
		return nil, err
	}
	p.startQueueProcessor()

	return p, nil
}

// PlayEffect plays a sound effect (fx) through the plugin's audio system. It optionally takes
// one or more callback functions that are executed once the sound effect playback is initiated.
//
// Parameters:
// - fx: The sound effect to be played, implementing the files.SoundFx interface.
// - callback: Optional variadic functions that will be executed after the sound effect starts playing.
//
// Returns:
// - An error if opening the sound effect fails, otherwise nil.
//
// This method is typically used for short, non-looping sound effects that do not interrupt, but overlap
// with the currently playing music or other sound effects.
func (p *Sound) PlayEffect(streams ...beep.Streamer) {
	if p.effectCount.Load() <= 5 {
		seq := make([]beep.Streamer, len(streams))
		for i, stream := range streams {
			seq[i] = beep.Seq(stream, beep.Callback(func() {
				p.effectCount.Add(-1)
			}))
		}
		p.effectCount.Add(1)
		p.play(seq...)
	} else {
		log.WithField("effects", p.effectCount.Load()).Warn("to many effect playing, skipping effect.")
	}

}

// PlayMusic plays a piece of music defined by the fx parameter, replacing any currently playing music.
// If there is music already playing, it fades out the current track before starting the new one.
// Callback functions can be specified to run after the new music stops playing.
//
// Parameters:
// - fx: The music track to be played, implementing the files.SoundFx interface.
// - callback: Optional variadic functions that will be executed after the music stops playing.
//
// Returns:
// - An error if opening the music track fails, otherwise nil (currently, errors during opening are ignored).
//
// This method manages the transition between tracks, ensuring that only one piece of music plays at a time.
// It introduces a fade-out effect for the currently playing music and a fade-in effect for the new track
// for a smooth transition. Additionally, the music is initially muted, played, and then faded in.
func (p *Sound) PlayMusic(stream *streamer.Playback, callback ...func()) error {
	if p.nowPlaying != nil {
		previous := p.nowPlaying
		go func() {
			time.Sleep(1 * time.Second)
			previous.FadeOut(2*time.Second, func() {
				previous.Stop()
			})
		}()
	}
	p.nowPlaying = stream
	p.play(beep.Seq(stream, beep.Callback(func() {
		stream.Close()
		for _, f := range callback {
			f()
		}

	})))
	return nil

}

func (p *Sound) StopMusic() {
	if p.nowPlaying != nil {
		p.nowPlaying.FadeOutAndStop(1 * time.Second)
	}
}

// Announce queues a spoken announcement generated from the provided paragraph text.
// The announcement is added to the plugin's queue for sequential playback.
//
// Parameters:
// - paragraph: The text to be converted into speech and announced.
//
// Returns:
// - An error if the announcement generation fails, otherwise nil.
//
// This method utilizes the plugin's announcer to convert text to speech and then enqueues
// the generated announcement for playback. It's typically used for narrating game events,
// instructions, or notifications.
func (p *Sound) Announce(paragraph announcer.Announcement, streamers ...beep.Streamer) error {
	allStreamers := make([]beep.Streamer, 0, len(streamers)+1)
	announce, err := p.announcer.Announce(paragraph)
	if err != nil {
		return err
	}

	allStreamers = append(allStreamers, announce)
	allStreamers = append(allStreamers, beep.Callback(func() {
		announce.Close()
	}))
	allStreamers = append(allStreamers, streamers...)

	p.queue <- beep.Seq(allStreamers...)
	return nil
}

// OptionalAnnouncement performs a spoken announcement like the Announce method but only if
// the announcement queue is currently empty. This ensures that the provided announcement does
// not interrupt or queue behind other announcements.
//
// Parameters:
// - paragraph: The text to be announced, pending the queue's current state.
//
// This method is designed for announcements that are optional or lower priority, where it's
// preferable to avoid interrupting ongoing announcements or adding to a backlog. Examples
// might include ambient game world details or non-critical updates.
func (p *Sound) OptionalAnnouncement(paragraph announcer.Announcement) error {
	if len(p.queue) == 0 {
		return p.Announce(paragraph)
	}
	return nil
}

func initCommentaryEngine(config *Config, p *Sound) error {
	var err error
	if config.Plugin.Sound.Announcements.Enabled {
		switch config.Plugin.Sound.Announcements.Engine {
		case "playht":
			p.announcer, err = playht.New(p.config.Plugin.Sound.Announcements.PlayHT)
		case "elevenlabs":
			p.announcer, err = elevenlabs.New(p.config.Plugin.Sound.Announcements.ElevenLabs)
		default:
			return errors.New("Unknown announcements engine: " + p.config.Plugin.Sound.Announcements.Engine)
		}
	}

	return err
}

func (p *Sound) initSpeaker() error {
	return speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
}

func (p *Sound) startQueueProcessor() {
	go func() {
		var lock sync.Mutex
		for {
			select {
			case sound := <-p.queue:
				lock.Lock()
				if p.nowPlaying != nil {
					p.nowPlaying.Fade(-2, 500*time.Millisecond)
				}
				p.play(beep.Seq(sound, beep.Callback(func() {
					if p.nowPlaying != nil {
						p.nowPlaying.Fade(0, 500*time.Millisecond)
					}
					lock.Unlock()
				})))
			case <-p.context.Done():
				return
			}
		}
	}()
}

func (p *Sound) play(streams ...beep.Streamer) {
	speaker.Play(streams...)
}
