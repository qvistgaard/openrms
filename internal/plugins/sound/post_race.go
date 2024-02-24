package sound

import (
	"github.com/gopxl/beep"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/sounds"
	"math/rand"
	"time"
)

func (p *Plugin) postRaceSequence() {
	p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
		Fs:       announcements,
		Filename: "announcements/finished.txt",
		Random:   true,
	}, beep.Callback(func() {
		p.startPostRaceMusic()
		p.announceResults()
	}))
}

func (p *Plugin) announceResults() {
	leader := p.telemetry.Leader()

	if leader != nil {
		i := rand.Intn(2) + 2
		time.AfterFunc(time.Duration(i)*time.Second, func() {
			p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/result.txt",
				Random:   true,
				Data:     p.tracker.cars[leader.Get()].TemplateData(),
			}, beep.Callback(func() {
				p.announceFastestLap()
			}))
		})
	}
}

func (p *Plugin) announceFastestLap() {
	fastest := p.telemetry.FastestLap()
	if fastest != nil {
		i := rand.Intn(2) + 2
		time.AfterFunc(time.Duration(i)*time.Second, func() {
			p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/fastest_lap.txt",
				Random:   true,
				Data:     p.tracker.cars[fastest.Get()].TemplateData(),
			})
		})
	}
}

func (p *Plugin) startPostRaceMusic() *time.Timer {
	return time.AfterFunc(0, func() {
		win := sounds.PostRaceMusic()
		p.sound.PlayMusic(win)
	})
}
