package streamer

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
)

var SampleRate = beep.SampleRate(44100)

func CreateStreamer(streamer beep.StreamSeekCloser, sampleRate beep.SampleRate) *Playback {

	gain := &effects.Gain{
		Streamer: streamer, // beep.Resample(4, sampleRate, SampleRate, streamer),
		Gain:     0,
	}

	ctrl := &beep.Ctrl{
		Streamer: gain,
		Paused:   false,
	}

	return &Playback{
		sampleRate: sampleRate,
		fileStream: streamer,
		ctrl:       ctrl,
		gain:       gain,
		volume: &effects.Volume{
			Streamer: ctrl,
			Base:     2,
			Volume:   0,
			Silent:   false,
		},
	}
}
