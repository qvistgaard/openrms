package state

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/telemetry"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	CarPitLaneSpeed          = "car-pit-lane-speed"
	CarMinSpeed              = "car-min-speed"
	CarMaxSpeed              = "car-max-speed"
	CarConfigMaxSpeed        = "car-config-max-speed"
	CarConfigMinSpeed        = "car-config-min-speed"
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
	maxSpeed := Speed(100)
	if v, ok := settings["max-speed"].(int); ok {
		maxSpeed = Speed(v)
	} else {
		log.WithField("value", settings["max-speed"]).
			WithField("type", fmt.Sprintf("%T", settings["max-speed"])).
			Warn("failed to set max-speed, incorrect type; must be integer")
	}
	c.Create(CarConfigMaxSpeed, maxSpeed)
	c.Create(CarMaxSpeed, maxSpeed)

	minSpeed := Speed(0)
	if v, ok := settings["min-speed"].(int); ok {
		minSpeed = Speed(v)
	} else {
		log.Warn("failed to set min-speed, incorrect type; must be integer")
	}
	c.Create(CarConfigMinSpeed, minSpeed)
	c.Create(CarMinSpeed, minSpeed)

	pitMaxSpeed := Speed(100)
	if p, ok := settings["pit"].(map[interface{}]interface{}); ok {
		if v, ok := p["max-speed"].(int); ok {
			pitMaxSpeed = Speed(v)
		}
	} else {
		log.Warn("failed to set max-pit-speed, incorrect type; must be integer")
	}
	c.Create(CarPitLaneSpeed, pitMaxSpeed)

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
