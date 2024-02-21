package sounds

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestLap(t *testing.T) {
	t.SkipNow()
	complete := make(chan bool, 1)

	speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
	log.Info("Start Playing")
	speaker.Play(beep.Seq(Lap(), beep.Callback(func() {
		complete <- true
	})))

	<-complete
}

func TestHeroic(t *testing.T) {
	t.SkipNow()

	complete := make(chan bool, 1)

	speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
	drivingToWin := DrivingToWin()

	log.Info(drivingToWin.Len())
	speaker.Play(beep.Seq(drivingToWin, beep.Callback(func() {
		complete <- true
	})))

	<-complete
}
