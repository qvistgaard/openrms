package sounds

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	"testing"
	"time"
)

func TestBeeps(t *testing.T) {
	t.SkipNow()
	complete := make(chan bool, 1)

	speaker.Init(streamer.SampleRate, streamer.SampleRate.N(time.Second/20))
	speaker.Play(beep.Seq(Beeps(), beep.Callback(func() {
		complete <- true
	})))

	<-complete
}
