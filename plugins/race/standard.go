package race

import (
	"github.com/google/uuid"
)

type StandardRace struct {
}

func (r StandardRace) id() uuid.UUID {
	return uuid.New()
}
