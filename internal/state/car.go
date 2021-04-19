package state

import (
	"time"
)

const (
	CarPitLaneSpeed          = "car-pit-lane-speed"
	CarMinSpeed              = "car-min-speed"
	CarMaxSpeed              = "car-max-speed"
	CarConfigMaxSpeed        = "car-config-max-speed"
	CarMaxBreaking           = "car-max-breaking"
	CarOnTrack               = "car-ontrack"
	CarControllerLink        = "car-controller-link"
	CarLapNumber             = "car-lap-number"
	CarLapTime               = "car-lap-time"
	CarInPit                 = "car-in-pit"
	CarReset                 = "car-reset"
	CarEventSequence         = "car-event-sequence"
	ControllerTriggerValue   = "controller-trigger-value"
	ControllerBtnUp          = "controller-btn-up"
	ControllerBtnDown        = "controller-btn-down"
	ControllerBtnTrackCall   = "controller-btn-track-call"
	ControllerBatteryWarning = "controller-battery-warning"
)

func CreateCar(race *Race, id uint8, settings map[string]interface{}, rules []Rule) *Car {
	c := new(Car)
	c.id = id
	c.race = race
	c.settings = settings
	c.state = CreateInMemoryRepository()
	c.Create(CarEventSequence, uint(0))
	c.Create(CarConfigMaxSpeed, uint8(255))
	for _, r := range rules {
		r.InitializeCarState(c)
	}
	for _, s := range c.state.All() {
		s.initialize()
	}
	return c
}

type Car struct {
	id       uint8
	settings map[string]interface{}
	state    Repository
	race     *Race
}

type CarChanges struct {
	Car     uint8     `json:"id"`
	Changes []Change  `json:"changes"`
	Time    time.Time `json:"time"`
}

func (c *Car) Race() *Race {
	return c.race
}

func (c *Car) ResetStateChangeStatus() {
	c.state.ResetChanges()
}

func (c *Car) Changes() CarChanges {
	changes := CarChanges{
		Car:     c.id,
		Changes: make([]Change, 0),
		Time:    time.Now(),
	}
	for k, v := range c.state.Changes() {
		changes.Changes = append(changes.Changes, Change{
			Name:  k,
			Value: v.Get(),
		})
	}
	return changes
}

func (c *Car) Get(state string) interface{} {
	return c.state.Get(state).Get()
}
func (c *Car) Set(state string, value interface{}) {
	c.state.Get(state).Set(value)
}
func (c *Car) Create(state string, value interface{}) {
	c.state.Create(state, value)
}

func (c *Car) SetDefault(state string) {
	get := c.state.Get(state)
	get.Set(get.Initial())
}

func (c *Car) Id() uint8 {
	return c.id
}

func (c *Car) Subscribe(state string, s Subscriber) {
	c.state.Get(state).Subscribe(s)
}
