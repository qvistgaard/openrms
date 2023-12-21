package playht

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Engine struct {
	format string
	voice  *Voice
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

func New() *Engine {
	engine := &Engine{}

	voice, _ := engine.getVoice("Oliver (Advertising)")
	engine.voice = voice

	return engine
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
	filename := "voices.json"

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
		return nil, errors.WithMessage(err, "Unable to load voices")
	}

	// Print the response
	for _, voice := range voices {
		fmt.Printf("%+v\n", voice)
		if strings.Compare(voice.Name, name) == 0 {
			log.WithField("voice", voice.Name).Info("Speaker found")
			return &voice, nil
		}
	}
	return nil, errors.New("Could not find voice")
}

/**
We are underway, here at <track name>!
there's contact there's all sorts of pandemonium
oh big accident of the <turn
oh contact that was a big a hit,look at the damage at the front of car number <number>
Car number <number> is almost out of fuel.
Oh no. car number 99 has run out of fuel
oh oh, that is damage big time for car number 77
Car number <number> has gone of the track
The lights go green, and we are under way here at <track name>
The lights go green, Blast off here at <track name>
Car number <number> takes the lead on lap <lap number>
And across the line! it's a win for <team name> in car number <number>
It's a win for <team name> and car number <number>
Oh my goodness, the wheels actually has come of the car number <number>
*/

func (e *Engine) downloadSpeak(speak string) (*os.File, error) {
	filename := fmt.Sprintf("%x.mp3", md5.Sum([]byte(speak+e.voice.Manifest)))
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
		log.Info("Speak exists, using cached version")
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
