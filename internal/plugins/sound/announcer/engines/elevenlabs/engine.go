package elevenlabs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Engine struct {
	voice  *Voice
	config *ElevenLabsConfig
	apiKey string
}

type Voices struct {
	Voice []Voice `json:"voices"`
}

type Voice struct {
	VoiceId string `json:"voice_id"`
	Name    string `json:"name"`
}

func New(config *ElevenLabsConfig) (*Engine, error) {
	engine := &Engine{
		apiKey: config.ApiKey,
		config: config,
	}
	newpath := filepath.Join(".", config.Cache)
	err := os.MkdirAll(newpath, os.ModePerm)

	voice, err := engine.getVoice(config.Voice)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to get voice")
	}
	engine.voice = voice
	return engine, nil
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

type GenerateSpeak struct {
	Text     string        `json:"text"`
	Settings VoiceSettings `json:"voice_settings"`
}

func (e *Engine) getVoice(name string) (*Voice, error) {
	var voices = &Voices{}
	filename := e.config.Cache + "/eleven-labs-voices.json"

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {

		log.Info("Downloading speaker information")
		resp, err := grequests.Get("https://api.elevenlabs.io/v1/voices", &grequests.RequestOptions{
			Headers: map[string]string{
				"xi-api-key": e.apiKey,
			},
		})

		if err != nil {
			return nil, errors.WithMessage(err, "Unable to load voices")
		}
		if !resp.Ok {
			return nil, errors.New("Failed to load voices: " + resp.String())
		}

		err = resp.DownloadToFile(filename)

		if err != nil {
			return nil, errors.WithMessage(err, "Unable to load speak")
		}
	}

	open, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open speak")

	}

	err = json.NewDecoder(open).Decode(&voices)
	if err != nil {
		dir, err := os.Getwd()
		return nil, errors.WithMessage(err, "Unable to load voices: "+dir+filename)
	}

	// Print the response
	for _, voice := range voices.Voice {
		// fmt.Printf("%+v\n", voice)
		if strings.Compare(voice.Name, name) == 0 {
			log.WithField("voice", voice.Name).Info("Speaker found")
			return &voice, nil
		}
	}
	return nil, errors.New("Could not find voice")
}

func (e *Engine) downloadSpeak(speak string) (*os.File, error) {
	filename := fmt.Sprintf("%s/%x.mp3", e.config.Cache, md5.Sum([]byte(speak+e.voice.VoiceId)))
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		log.Info("Speak not found. Generating new speak")

		resp, err := grequests.Post("https://api.elevenlabs.io/v1/text-to-speech/"+e.voice.VoiceId, &grequests.RequestOptions{
			Headers: map[string]string{
				"xi-api-key":   e.apiKey,
				"Accept":       "audio/mpeg",
				"Content-Type": "application/json",
			},
			JSON: GenerateSpeak{
				Text: speak,
				Settings: VoiceSettings{
					Stability:       0.5,
					SimilarityBoost: 0.5,
				},
			},
		})

		if !resp.Ok {
			return nil, errors.New("Failed to load speak: " + resp.String())
		}

		err = resp.DownloadToFile(filename)

		if err != nil {
			return nil, errors.WithMessage(err, "Unable to load speak")
		}

		log.WithField("filename", filename).WithField("speak", speak).Info("Speak successfully created")
	} else {
		log.WithField("filename", filename).WithField("speak", speak).Info("Speak exists, using cached version")
	}

	open, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open speak")

	}

	return open, nil

}

func (e *Engine) Announce(announcement announcer.Announcement) (*streamer.Playback, error) {
	downloadSpeak, err := e.downloadSpeak(announcement.Get())

	fromFs, err := streamer.LoadMp3FromFsFile(downloadSpeak)
	if err != nil {
		fromFs.Gain(1.5)
	}
	return fromFs, err

}
