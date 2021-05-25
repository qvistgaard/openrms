package state

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/telemetry"
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
	CarLap                   = "car-lap"
	CarInPit                 = "car-in-pit"
	CarReset                 = "car-reset"
	CarEventSequence         = "car-event-sequence"
	ControllerTriggerValue   = "controller-trigger-value"
	ControllerBtnUp          = "controller-btn-up"
	ControllerBtnDown        = "controller-btn-down"
	ControllerBtnTrackCall   = "controller-btn-track-call"
	ControllerBatteryWarning = "controller-battery-warning"
)

type CarId uint8
type Speed uint8
type TriggerValue uint8
type LapNumber uint
type RaceTimer time.Duration
type LapTime time.Duration
type Lap struct {
	LapNumber LapNumber `json:"lap-number"`
	RaceTimer RaceTimer `json:"race-timer"`
	LapTime   LapTime   `json:"lap-time"`
}

func CreateCar(id CarId, settings map[string]interface{}, rules Rules) *Car {
	c := new(Car)
	c.id = id
	c.settings = settings
	c.state = CreateInMemoryRepository(c)
	c.Create(CarEventSequence, uint(0))
	c.Create(CarOnTrack, false)
	for _, r := range rules.All() {
		r.InitializeCarState(c)
	}
	for _, s := range c.state.All() {
		s.initialize()
	}
	// TODO: make these values configurable
	c.Create(CarConfigMaxSpeed, Speed(255))
	c.Create(CarMaxSpeed, Speed(255))
	c.Create(CarPitLaneSpeed, Speed(75))
	return c
}

type Car struct {
	id       CarId
	settings map[string]interface{}
	state    Repository
}

type CarState struct {
	Car     CarId     `json:"id"`
	Changes []Change  `json:"changes"`
	Time    time.Time `json:"time"`
}

func (c *Car) Settings(v interface{}) error {
	return mapstructure.Decode(c.settings, v)
}

func (c *Car) ResetStateChangeStatus() {
	c.state.ResetChanges()
}

func (c *Car) State() CarState {
	return c.mapState(c.state.All())

}

func (c *Car) Changes() CarState {
	return c.mapState(c.state.Changes())
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

func (c *Car) Id() CarId {
	return c.id
}

func (c *Car) Subscribe(state string, s Subscriber) {
	c.state.Get(state).Subscribe(s)
}

func (c *Car) mapState(state map[string]StateInterface) CarState {
	changes := CarState{
		Car:     c.id,
		Changes: make([]Change, 0),
		Time:    time.Now(),
	}
	for k, v := range state {
		changes.Changes = append(changes.Changes, Change{
			Name:  k,
			Value: v.Get(),
		})
	}
	return changes
}

func (l *Lap) Metrics() []telemetry.Metric {
	return []telemetry.Metric{
		{Name: "car-lap-time", Value: l.LapTime},
		{Name: "car-lap-race-timer", Value: l.RaceTimer},
		{Name: "car-lap-number", Value: l.LapNumber},
	}
}
