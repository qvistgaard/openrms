package utils

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"io/fs"
)

func PlayAudioFile(file fs.File, callback func()) error {
	sampleRate := beep.SampleRate(44100)
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return err
	}
	resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		streamer.Close()
		callback()
	})))

	return nil
}
