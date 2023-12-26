package names

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/utils"
	log "github.com/sirupsen/logrus"
)

//go:embed drivers.txt
//go:embed teams.txt
//go:embed manufactures.txt
var names embed.FS

func RandomDriver() string {
	name, err := utils.RandomLine(names, "drivers.txt")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}

func RandomTeam() string {
	name, err := utils.RandomLine(names, "teams.txt")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}

func RandomColor() string {
	name, err := utils.RandomLine(names, "colors.txt")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}

func RandomManufacture() string {
	name, err := utils.RandomLine(names, "manufactures.txt")
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}
