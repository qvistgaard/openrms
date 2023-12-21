package playht

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func Test_getVoice(t *testing.T) {
	engine := New()
	stream, _ := engine.Announce("We are underway, here at Monza!")

	streamer, format, err := mp3.Decode(stream)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	eq := effects.NewEqualizer(streamer, beep.SampleRate(44100), effects.MonoEqualizerSections{

		// {F0: 40, Bf: 0, GB: 0, G0: 0, G: -100},
		// {F0: 60, Bf: 10000, GB: 0, G0: 0, G: 0},
		/*		{F0: 250, Bf: 5, GB: 3, G0: 0, G: 10},
				{F0: 300, Bf: 5, GB: 3, G0: 0, G: 12},
				{F0: 350, Bf: 5, GB: 3, G0: 0, G: 14},*/
		{F0: 100, Bf: 200, GB: 5, G0: 5, G: -100},
		{F0: 2500, Bf: 500, GB: 1, G0: 1, G: 0},
		{F0: 5000, Bf: 10, GB: 12, G0: 12, G: -100},
		// {F0: 6000, Bf: 50, GB: 3, G0: 0, G: -100},
	})

	// Initialize the speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Play the file
	done := make(chan bool)
	speaker.Play(beep.Seq(eq, beep.Callback(func() {
		done <- true
	})))

	// Wait for the playback to finish
	<-done
}
