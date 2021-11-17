package state

/*
import (
	"time"
)

const (
	RaceStatus                  = "race-status"
	RaceTime                    = "race-time"
	CourseMaxSpeed              = "course-max-speed"
	CourseLength                = "course-length"
	RMSStatus                   = "rms-status"
	Initialized                 = "initialized"
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
		MaxSpeed uint8 `mapstructure:"max-speed"`
		Length   int
		MaxSpeed  struct {
			LapCounting struct {
				Enabled bool `mapstructure:"enabled"`
				OnEntry bool `mapstructure:"on-entry"`
			} `mapstructure:"lap-counting"`
		}
	}
}

func CreateCourse(config *CourseConfig) *Course {
	course := new(Course)
	course.state = CreateInMemoryRepository(course)

	course.state.Create(PitlaneLapCounting, config.Course.MaxSpeed.LapCounting.Enabled)
	course.state.Create(PitlaneLapCountingOnEntry, config.Course.MaxSpeed.LapCounting.OnEntry)
	course.state.Create(CourseLength, config.Course.Length)
	course.state.Create(CourseMaxSpeed, Speed(config.Course.MaxSpeed))
	course.state.Create(RaceStatus, RaceStatusStopped)

	course.state.SetDefaults()
	for _, r := range rules.All() {
		r.InitializeCourseState(course)
	}
	for _, s := range course.state.All() {
		s.initialize()
	}
	return course
}

type CourseState struct {
	Changes []Change  `json:"changes"`
	time    time.time `json:"time"`
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

func (c *Course) Changes() CourseState {
	return c.mapState(c.state.Changes())

}

func (c *Course) Subscribe(state string, s Subscriber) {
	c.state.Get(state).Subscribe(s)
}

func (c *Course) IsChanged(state string) bool {
	return c.state.Get(state).Changed()
}

func (c *Course) State() CourseState {
	return c.mapState(c.state.All())
}

func (c *Course) mapState(state map[string]StateInterface) CourseState {
	changes := CourseState{
		Changes: []Change{},
		time:    time.Now(),
	}
	for k, v := range state {
		changes.Changes = append(changes.Changes, Change{
			Name:  k,
			Value: v.Get(),
		})
	}
	return changes
}

type Settings struct {
}
*/
