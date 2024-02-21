package sounds

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/plugins/sound/streamer"
	"math/rand"
	"time"
)

//go:embed files/race-care-151963.mp3
//go:embed files/driving-to-win-16372.mp3
//go:embed files/action-stylish-rock-dedication-15038.mp3
//go:embed files/cool-sport-rock-beat-95429.mp3
//go:embed files/energetic-indie-rock-jump-112179.mp3
//go:embed files/iron-man-190508.mp3
//go:embed files/powerful-stylish-stomp-rock-lets-go-114255.mp3
//go:embed files/inspiring-epic-motivation-cinematic-trailer-11218.mp3
//go:embed files/cinematic-epic-trailer-background-music-123922.mp3
//go:embed files/powerful-victory-trailer-103656.mp3
//go:embed files/the-epic-inspiring-153397.mp3
//go:embed files/the-shield-111353.mp3
var files embed.FS

func Lap() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/race-care-151963.mp3")
	fs.SeekToPositionInDuration(900 * time.Millisecond)
	fs.Gain(-.5)
	return fs
}

func DrivingToWin() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/driving-to-win-16372.mp3")
	return fs
}

func ActionStylishRock() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/action-stylish-rock-dedication-15038.mp3")
	return fs
}

func CoolSportsRockBeat() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/cool-sport-rock-beat-95429.mp3")
	return fs
}

func EnergeticIndiRockJump() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/energetic-indie-rock-jump-112179.mp3")
	return fs
}

func IronMan() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/iron-man-190508.mp3")
	return fs
}

func StompRock() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/powerful-stylish-stomp-rock-lets-go-114255.mp3")
	return fs
}

func InspiringEpic() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/inspiring-epic-motivation-cinematic-trailer-11218.mp3")
	fs.SoftLen((1 * time.Minute) + (40 * time.Second) + (800 * time.Millisecond))
	return fs
}

func EpicTrailer() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/cinematic-epic-trailer-background-music-123922.mp3")
	fs.SoftLen((1 * time.Minute) + (6 * time.Second) + (400 * time.Millisecond))
	return fs
}

func PowerfulVictory() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/powerful-victory-trailer-103656.mp3")
	fs.SoftLen((2 * time.Minute) + (1 * time.Second) + (800 * time.Millisecond))
	return fs
}

func TheEpic() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/the-epic-inspiring-153397.mp3")
	fs.SoftLen((2 * time.Minute) + (0 * time.Second) + (400 * time.Millisecond))
	return fs
}

func TheShield() *streamer.Playback {
	fs, _ := streamer.LoadMp3FromFs(files, "files/the-shield-111353.mp3")
	fs.SoftLen((2 * time.Minute) + (0 * time.Second) + (400 * time.Millisecond))
	return fs
}

func PostRaceMusic() *streamer.Playback {
	songs := [...]func() *streamer.Playback{
		DrivingToWin, ActionStylishRock, CoolSportsRockBeat, EnergeticIndiRockJump,
		IronMan, StompRock,
	}
	line := rand.Intn(len(songs) - 1)
	return songs[line]()

}

func EpicRise() *streamer.Playback {
	songs := [...]func() *streamer.Playback{
		InspiringEpic, EpicTrailer, PowerfulVictory, TheEpic, TheShield,
	}
	line := rand.Intn(len(songs) - 1)
	return songs[line]()

}
