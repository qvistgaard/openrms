package state

import (
	"gopkg.in/yaml.v2"
	"time"
)

const (
	RaceStatus                  = "race-status"
	CourseMaxSpeed              = "course-length-max-speed"
	CourseLength                = "course-length"
	PitlaneLapCounting          = "cource-pitlane-lap-counting"
	PitlaneLapCountingOnEntry   = "cource-pitlane-lap-counting-on-entry"
	RaceStatusStopped           = uint8(0x00)
	RaceStatusPaused            = uint8(0x01)
	RaceStatusRunning           = uint8(0x02)
	RaceStatusFlaggedLCDisabled = uint8(0x04)
	RaceStatusFlaggedLCEnabled  = uint8(0x08)
)

type CourseConfig struct {
	Course struct {
		MaxSpeed uint8 `yaml:"max-speed"`
		Length   int
		PitLane  struct {
			LapCounting struct {
				Enabled bool
				Entry   bool
			} `yaml:"lap-counting"`
		}
	}
}

func CreateCourseFromConfig(config []byte, rules Rules) (*Course, error) {
	c := &CourseConfig{}
	perr := yaml.Unmarshal(config, c)
	if perr != nil {
		return nil, perr
	}
	return CreateCourse(c, rules), nil
}

func CreateCourse(config *CourseConfig, rules Rules) *Course {
	course := new(Course)
	course.state = CreateInMemoryRepository(course)

	course.state.Create(PitlaneLapCounting, config.Course.PitLane.LapCounting.Enabled)
	course.state.Create(PitlaneLapCountingOnEntry, config.Course.PitLane.LapCounting.Entry)
	course.state.Create(CourseLength, config.Course.Length)
	course.state.Create(CourseMaxSpeed, config.Course.MaxSpeed)

	course.state.SetDefaults()
	for _, r := range rules.All() {
		r.InitializeCourseState(course)
	}
	for _, s := range course.state.All() {
		s.initialize()
	}
	return course
}

type CourseChanges struct {
	Changes []Change  `json:"changes"`
	Time    time.Time `json:"time"`
}

type Course struct {
	settings map[string]interface{}
	state    Repository
}

func (c *Course) Get(state string) interface{} {
	return c.state.Get(state).Get()
}
func (c *Course) Set(state string, value interface{}) {
	c.state.Get(state).Set(value)
}

func (c *Course) ResetStateChangeStatus() {
	c.state.ResetChanges()
}

func (c *Course) Changes() CourseChanges {
	stateChanges := c.state.Changes()
	changes := CourseChanges{
		Changes: []Change{},
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

func (c *Course) Subscribe(state string, s Subscriber) {
	c.state.Get(state).Subscribe(s)
}

type Settings struct {
}