package streamer

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"time"
)

type Playback struct {
	sampleRate beep.SampleRate
	softLen    *time.Duration
	ctrl       *beep.Ctrl
	volume     *effects.Volume
	gain       *effects.Gain
	playing    bool
	fileStream beep.StreamSeekCloser
	cancelFade bool
	stopping   bool
}

func (p *Playback) Stream(samples [][2]float64) (n int, ok bool) {
	return p.volume.Stream(samples)
}

func (p *Playback) Err() error {
	return p.volume.Err()
}

func (p *Playback) Len() int {
	return p.fileStream.Len()
}

func (p *Playback) Position() int {
	return p.fileStream.Position()
}

func (p *Playback) Seek(position int) error {
	return p.fileStream.Seek(position)
}

func (p *Playback) Close() error {
	return p.fileStream.Close()
}

func (p *Playback) Pause() {
	p.ctrl.Paused = true
}

func (p *Playback) Gain(gain float64) {
	p.gain.Gain = gain
}

func (p *Playback) SoftLen(len time.Duration) {
	p.softLen = &len
}

func (p *Playback) FadeOut(duration time.Duration, callback ...func()) {
	if !p.stopping {
		p.Fade(-6.0, duration, func() {
			p.volume.Silent = true
			executeCallbacks(callback)
		})
	} else {
		executeCallbacks(callback)
	}
}

func (p *Playback) FadeIn(duration time.Duration, callback ...func()) {
	if !p.stopping {
		p.volume.Volume = -6.0
		p.volume.Silent = false
		p.Fade(0, duration, callback...)
	} else {
		executeCallbacks(callback)
	}
}

func (p *Playback) Fade(targetVolume float64, duration time.Duration, callback ...func()) {
	p.cancelFade = true
	go func() {
		p.cancelFade = false
		delay := 20 * time.Millisecond
		steps := int(duration / delay)
		startVolume := p.volume.Volume
		stepSize := (targetVolume - startVolume) / float64(steps)

		for i := 0; i <= steps && !p.cancelFade; i++ {
			p.volume.Volume = startVolume + stepSize*float64(i)
			time.Sleep(delay)
		}
		executeCallbacks(callback)
	}()
}

func executeCallbacks(callback []func()) {
	for _, f := range callback {
		f()
	}
}

func (p *Playback) Mute() {
	p.volume.Silent = true
}

func (p *Playback) Unmute() {
	p.volume.Silent = false
}

func (p *Playback) PositionInSeconds() time.Duration {
	return p.sampleRate.D(p.Position()).Round(time.Second)
}

func (p *Playback) SeekToPositionInDuration(duration time.Duration) error {
	return p.Seek(p.sampleRate.N(duration))
}

func (p *Playback) SoftLenAsDuration() time.Duration {
	if p.softLen != nil {
		return *p.softLen
	}
	return p.LenAsDuration()
}

func (p *Playback) LenAsDuration() time.Duration {
	return p.sampleRate.D(p.fileStream.Len())
}

func (p *Playback) FadeOutAndStop(duration time.Duration) {
	p.stopping = true
	p.FadeOut(duration, func() {
		p.Stop()
	})
}

func (p *Playback) Stop() {
	p.ctrl.Streamer = nil
}

func (p *Playback) IsPlaying() bool {
	return !p.ctrl.Paused
}

func (p *Playback) IsPaused() bool {
	return p.ctrl.Paused
}

func (p *Playback) SampleRate() beep.SampleRate {
	return p.sampleRate
}
