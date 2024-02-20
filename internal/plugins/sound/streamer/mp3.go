package streamer

import (
	"github.com/gopxl/beep/mp3"
	"io/fs"
)

func LoadMp3FromFs(fs fs.FS, filename string) (*Playback, error) {
	file, _ := fs.Open(filename)
	return LoadMp3FromFsFile(file)
}

func LoadMp3FromFsFile(file fs.File) (*Playback, error) {
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return nil, err
	}
	return CreateStreamer(streamer, format.SampleRate), nil
}

/*
type Mp3File struct {
	FileDuration time.Duration
	File         fs.File
	Gain         float64
}

func (f Mp3File) Duration() time.Duration {
	return f.FileDuration
}

func (f Mp3File) Streamer() *Playback {
	streamer, format, err := mp3.Decode(f.File)
	if err != nil {
		return nil
	}

	resampled := beep.Resample(4, format.SampleRate, beep.SampleRate(44100), streamer)

	ctrl := &beep.Ctrl{
		Streamer: &effects.Gain{
			Streamer: resampled,
			Gain:     f.Gain,
		},
		Paused: false,
	}
	return &Playback{
		sampleRate: format.SampleRate,
		duration:   f.FileDuration,
		fileStream: streamer,
		ctrl:       ctrl,
		volume: &effects.Volume{
			Streamer: ctrl,
			Base:     2,
			Volume:   0,
			Silent:   false,
		},
	}
}
*/
