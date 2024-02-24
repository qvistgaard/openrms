package sound

import (
	"github.com/gopxl/beep"
	"github.com/qvistgaard/openrms/internal/plugins/sound/announcer"
	"github.com/qvistgaard/openrms/internal/plugins/sound/sounds"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func (p *Plugin) postRaceSequence() {
	err := p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
		Fs:       announcements,
		Filename: "announcements/finished.txt",
		Random:   true,
	}, beep.Callback(func() {
		p.startPostRaceMusic()
		p.announceResults()
	}))
	if err != nil {
		log.Error(err)
	}
}

func (p *Plugin) announceResults() {
	leader := p.telemetry.Leader()

	if leader != nil {
		i := rand.Intn(2) + 2
		time.AfterFunc(time.Duration(i)*time.Second, func() {
			err := p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/result.txt",
				Random:   true,
				Data:     p.tracker.cars[leader.Get()].TemplateData(),
			}, beep.Callback(func() {
				p.announceFastestLap()
			}))
			if err != nil {
				log.Error(err)
			}
		})
	}
}

func (p *Plugin) announceFastestLap() {
	fastest := p.telemetry.FastestLap()
	if fastest != nil {
		i := rand.Intn(2) + 2
		time.AfterFunc(time.Duration(i)*time.Second, func() {
			err := p.sound.Announce(&announcer.ReadFileTemplateAnnouncement{
				Fs:       announcements,
				Filename: "announcements/fastest_lap.txt",
				Random:   true,
				Data:     p.tracker.cars[fastest.Get()].TemplateData(),
			})
			if err != nil {
				log.Error(err)
			}
		})
	}
}

func (p *Plugin) startPostRaceMusic() *time.Timer {
	return time.AfterFunc(0, func() {
		win := sounds.PostRaceMusic()
		err := p.sound.PlayMusic(win)
		if err != nil {
			log.Error(err)
		}
	})
}
