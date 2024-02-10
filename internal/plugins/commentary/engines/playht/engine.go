package playht

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Engine struct {
	voice  *Voice
	config *PlayHTConfig
	apiKey string
	userId string
}

type Voice struct {
	Name     string `json:"name"`
	Manifest string `json:"id"`
	Accent   string `json:accent`
	Age      string `json:age`
	Language string `json:language`
	Style    string `json:style`
}

func New(config *PlayHTConfig) (*Engine, error) {
	engine := &Engine{
		apiKey: config.ApiKey,
		userId: config.UserId,
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

type GenerateSpeak struct {
	Text       string  `json:"text"`
	Voice      string  `json:"voice"`
	Format     string  `json:"output_format"`
	Engine     string  `json:"voice_engine"`
	Speed      float32 `json:"speed"`
	SampleRate int32   `json:"sample_rate"`
}

func (e *Engine) getVoice(name string) (*Voice, error) {
	var voices = make([]Voice, 0)
	filename := e.config.Cache + "/voices.json"

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {

		log.Info("Downloading speaker information")
		resp, err := grequests.Get("https://api.play.ht/api/v2/voices", &grequests.RequestOptions{
			Headers: map[string]string{
				"AUTHORIZATION": e.apiKey,
				"X-USER-ID":     e.userId,
			},
		})
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to load voices")
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
	for _, voice := range voices {
		// fmt.Printf("%+v\n", voice)
		if strings.Compare(voice.Name, name) == 0 {
			log.WithField("voice", voice.Name).Info("Speaker found")
			return &voice, nil
		}
	}
	return nil, errors.New("Could not find voice")
}

func (e *Engine) downloadSpeak(speak string) (*os.File, error) {
	filename := fmt.Sprintf("%s/%x.mp3", e.config.Cache, md5.Sum([]byte(speak+e.voice.Manifest)))
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		log.Info("Speak not found. Generating new speak")

		resp, err := grequests.Post("https://api.play.ht/api/v2/tts/stream", &grequests.RequestOptions{
			Headers: map[string]string{
				"AUTHORIZATION": e.apiKey,
				"X-USER-ID":     e.userId,
			},
			JSON: GenerateSpeak{
				Text:       speak,
				Voice:      e.voice.Manifest,
				Format:     "mp3",
				Engine:     "PlayHT2.0",
				Speed:      float32(1.1),
				SampleRate: 48000,
			},
		})

		err = resp.DownloadToFile(filename)

		if err != nil {
			return nil, errors.WithMessage(err, "Unable to load speak")
		}

		log.WithField("filename", filename).Info("Speak successfully created")
	} else {
		log.WithField("filename", filename).Info("Speak exists, using cached version")
	}

	open, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open speak")

	}

	return open, nil

}

func (e *Engine) Announce(speak string) (*os.File, error) {
	return e.downloadSpeak(speak)
}
