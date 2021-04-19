package state

import "time"

const (
	RaceStatus                  = "race-status"
	RaceMaxSpeed                = "max-speed"
	RaceStatusStopped           = uint8(0x00)
	RaceStatusPaused            = uint8(0x01)
	RaceStatusRunning           = uint8(0x02)
	RaceStatusFlaggedLCDisabled = uint8(0x04)
	RaceStatusFlaggedLCEnabled  = uint8(0x08)
)

func CreateRace(settings map[string]interface{}, rules []Rule) *Race {
	race := new(Race)
	race.settings = settings
	race.state = CreateInMemoryRepository()
	race.state.SetDefaults()
	for _, r := range rules {
		r.InitializeRaceState(race)
	}
	for _, s := range race.state.All() {
		s.initialize()
	}
	return race
}

type RaceChanges struct {
	Changes []Change  `json:"changes"`
	Time    time.Time `json:"time"`
}

type Race struct {
	settings map[string]interface{}
	state    Repository
}

func (r *Race) Get(state string) interface{} {
	return r.state.Get(state).Get()
}
func (r *Race) Set(state string, value interface{}) {
	r.state.Get(state).Set(value)
}

func (r *Race) ResetStateChangeStatus() {
	r.state.ResetChanges()
}

func (r *Race) Changes() RaceChanges {
	stateChanges := r.state.Changes()
	changes := RaceChanges{
		Changes: make([]Change, len(stateChanges)),
		Time:    time.Now(),
	}
	for k, v := range stateChanges {
		changes.Changes = append(changes.Changes, Change{
			Name:  k,
			Value: v.Get(),
		})
	}
	return changes
}

type Settings struct {
}
