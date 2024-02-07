package names

import (
	"embed"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/utils"
	log "github.com/sirupsen/logrus"
)

//go:embed drivers.txt
//go:embed teams.txt
//go:embed manufactures.txt
var names embed.FS

func RandomDriver(id types.CarId) string {
	name, err := utils.GetLine(names, "drivers.txt", int(id))
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}

func RandomTeam(id types.CarId) string {
	name, err := utils.GetLine(names, "teams.txt", int(id))
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

func RandomManufacture(id types.CarId) string {
	name, err := utils.GetLine(names, "manufactures.txt", int(id))
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return name
}
