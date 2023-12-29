package commentary

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/commentary/engines/playht"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Plugin struct {
	config        *Config
	engine        Engine
	queue         chan string
	playerRunning sync.Mutex
	playerStarted bool
	sampleRate    beep.SampleRate
}

func New(config *Config) (*Plugin, error) {

	p := &Plugin{config: config}
	p.config = config
	p.queue = make(chan string, 20)
	p.sampleRate = beep.SampleRate(44100)

	var err error
	err = speaker.Init(p.sampleRate, p.sampleRate.N(time.Second/10))
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to create egine")
	}

	if config.Plugin.Commentary.Enabled {
		switch p.config.Plugin.Commentary.Engine {
		case "playht":
			p.engine, err = playht.New(p.config.Plugin.Commentary.PlayHT)
		default:
			return nil, errors.New("Unknown commentary engine: " + p.config.Plugin.Commentary.Engine)
		}
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Unable to create engine")
	}

	return p, nil
}

func (p *Plugin) Announce(paragraph string) {
	p.start()
	p.queue <- paragraph
}

func (p *Plugin) OptionalAnnouncement(paragraph string) {
	if len(p.queue) == 0 {
		p.Announce(paragraph)
	}
}

func (p *Plugin) start() {
	if !p.playerStarted {
		p.playerStarted = true
		go func() {
			defer func() { p.playerStarted = false }()
			for p.playerStarted {
				select {
				case paragraph := <-p.queue:
					if p.config.Plugin.Commentary.Enabled {
						err := p.play(paragraph)
						if err != nil {
							log.Error(err)
						}
					}
				case <-time.After(time.Duration(1) * time.Second):

				}
			}
		}()
	}
}

func (p *Plugin) play(paragraph string) error {
	p.playerRunning.Lock()
	log.WithField("paragraph", paragraph).Info("Playing announcement")
	stream, err := p.engine.Announce(paragraph)
	if err != nil {
		defer p.playerRunning.Unlock()
		return err
	}

	streamer, format, err := mp3.Decode(stream)
	if err != nil {
		defer p.playerRunning.Unlock()
		return err
	}
	resampled := beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		defer p.playerRunning.Unlock()
		streamer.Close()
	})))
	return nil
}

func (p *Plugin) Priority() int {
	return 10
}

func (p *Plugin) Name() string {
	return "commentary"
}
