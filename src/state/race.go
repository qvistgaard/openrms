package state

const (
	RaceStatus                  = "race-status"
	RaceMaxSpeed                = "max-speed"
	RaceStatusStopped           = uint8(0x00)
	RaceStatusPaused            = uint8(0x01)
	RaceStatusRunning           = uint8(0x02)
	RaceStatusFlaggedLCDisabled = uint8(0x04)
	RaceStatusFlaggedLCEnabled  = uint8(0x08)
)

func CreateRace(settings map[string]interface{}) *Race {
	r := new(Race)
	r.settings = settings
	r.state = CreateInMemoryRepository()
	return r
}

type Race struct {
	settings map[string]interface{}
	state    Repository
}

func (c *Race) State() Repository {
	return c.state
}

type Settings struct {
}
