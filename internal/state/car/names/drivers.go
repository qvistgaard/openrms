package names

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

//go:embed drivers.txt
//go:embed teams.txt
var names embed.FS

func RandomDriver() string {
	name, err := getRandomLine("report.gohtml")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}

func getRandomLine(file string) (string, error) {
	data, err := names.ReadFile(file)
	if err != nil {
		return "The Stig", err
	}
	split := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	intn := rand.Intn(len(split) - 1)

	return split[intn], nil
}

func RandomTeam() string {
	name, err := getRandomLine("teams.txt")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}
