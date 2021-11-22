package leaderboard

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

//go:embed drivers.txt
var content embed.FS

func getRandomDriver() string {
	file, err := content.ReadFile("drivers.txt")
	if err != nil {
		log.Error(err)
		return "The Stig"
	}
	split := strings.Split(strings.ReplaceAll(string(file), "\r\n", "\n"), "\n")
	intn := rand.Intn(len(split) - 1)

	return split[intn]
}
