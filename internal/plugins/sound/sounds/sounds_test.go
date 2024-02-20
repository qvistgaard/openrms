package sounds

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestBeeps(t *testing.T) {
	complete := make(chan bool, 1)

	speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
	speaker.Play(beep.Seq(Lap(), beep.Callback(func() {
		complete <- true
	})))

	<-complete
}

func TestHeroic(t *testing.T) {
	complete := make(chan bool, 1)

	speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
	drivingToWin := DrivingToWin()

	log.Info(drivingToWin.Len())

	// t.SkipNow()
	speaker.Play(beep.Seq(drivingToWin, beep.Callback(func() {
		complete <- true
	})))

	<-complete
}
