package sound

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

//go:embed files/ghosts.mp3
//go:embed files/horror.mp3
var files embed.FS

func Horror() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/horror.mp3")
	fs.Gain(-.1)
	return fs
}

func Ghosts() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/ghosts.mp3")
	return fs
}

func TestPlugin_PlayMusic(t *testing.T) {
	t.SkipNow()
	complete := make(chan bool, 1)
	plugin := Plugin{}
	plugin.initSpeaker()

	plugin.PlayMusic(Horror())
	time.Sleep(10 * time.Second)
	log.Info("Fade")
	plugin.PlayMusic(Ghosts(), func() {
		complete <- true
	})

	<-complete
	log.Info("Complete")

}
