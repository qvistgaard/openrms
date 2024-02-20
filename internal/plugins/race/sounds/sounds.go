package sounds

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
)

//go:embed files/beeps.mp3
var files embed.FS

func Beeps() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/beeps.mp3")
	fs.Gain(-.1)
	return fs
}
