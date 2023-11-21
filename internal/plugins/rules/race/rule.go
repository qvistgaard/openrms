package race

type Rule struct {
}

/*
import (
	"github.com/qvistgaard/openrms/internal/state"
	"github.com/qvistgaard/openrms/internal/telemetry"
	log "github.com/sirupsen/logrus"
	"time"
)

type Laps uint

type StageConfig struct {
	Type     *string
	Laps     *Laps
	Duration *string
	duration *time.Duration
}

type Rule struct {
	course *state.Course
	ready  map[state.CarId]*state.Car
	stages []StageConfig
}

type RaceConfirmationCounter struct {
	cars      uint8
	confirmed uint8
}

const (
	Countdown = "race-countdown"

	Cars      = "race-cars"
	CarsReady = "race-cars-ready"

	None = "none"

	ActiveView     = "race-state"
	Started   = "started"
	Running   = "running"
	Stopped   = "stopped"
	Paused    = "paused"
	TrackCall = "track-call"

	Confirmation = "race-confirmation"
	Unconfirmed  = "unconfirmed"
	Confirmed    = "confirmed"

	Ready = "car-ready"

	Mode       = "race-mode"
	Training   = "training"
	Qualifying = "qualifying"
	Race       = "race"

	Stage         = "race-stage"
	StageNumber   = "race-stage-number"
	StageDuration = "race-stage-duration"
	StageLaps     = "race-stage-laps"

	CountdownSteps = uint8(5)
)

func (r *Rule) Notify(v *state.Value) {
	if c, ok := v.Owner().(*state.Car); ok {
		r.handleReady(c, v)

		if s, vok := v.Get().(state.Lap); vok && v.Name() == state.CarLap {
			if d, dok := r.course.Get(StageLaps).(Laps); dok && r.course.Get(ActiveView) == Running {
				if d < Laps(s.LapNumber) {
					r.course.Set(ActiveView, Stopped)
				}
			}
		}
	}

	if c, ok := v.Owner().(*state.Course); ok {
		// Initialize Race when RMS is started
		if s, vok := v.Get().(string); vok && v.Name() == state.RMSStatus {
			if s == state.Initialized {
				c.Set(state.RaceStatus, state.RaceStatusRunning)
				log.WithField("race-state", c.Get(state.RaceStatus)).
					Info("race system initialized")
			}
		}

		if s, vok := v.Get().(uint8); vok && v.Name() == CarsReady {
			if s == r.course.Get(Cars) {
				log.WithField(state.RaceStatus, r.course.Get(state.RaceStatus)).
					Info("all cars confirmed")
				for cr := range r.ready {
					r.ready[cr].Set(Ready, false)
				}
				c.Set(CarsReady, uint8(0))
				c.Set(Confirmation, Confirmed)
			}
		}
		if s, vok := v.Get().(state.RaceTimer); vok && v.Name() == state.RaceTime {
			if d, dok := c.Get(StageDuration).(time.Duration); dok && c.Get(ActiveView) == Running {
				if d < time.Duration(s) {
					log.WithField(state.RaceStatus, r.course.Get(state.RaceStatus)).
						WithField(state.RaceTime, time.Duration(s)).
						Info("stage duration have been reached")
					c.Set(ActiveView, Stopped)
				}
			}
		}

		if s, vok := v.Get().(uint8); vok && v.Name() == StageNumber {
			if len(r.stages) >= int(s) {
				c.Set(Stage, r.stages[s])
			}
		}

		r.handleRaceState(c, v)
		r.handleCountdown(c, v)
	}

}

func (r *Rule) handleRaceState(c *state.Course, v *state.Value) {
	if s, ok := v.Get().(string); ok && v.Name() == ActiveView {
		if s == Started {
			if c.Get(Mode) == None {
				stage := c.Get(Stage).(StageConfig)
				c.Set(Mode, stage.Type)
				c.Set(StageDuration, nil)
				c.Set(StageLaps, nil)
				if stage.duration != nil {
					c.Set(StageDuration, *stage.duration)
				}
				if stage.Laps != nil {
					c.Set(StageLaps, *stage.Laps)
				}
			}
			c.Set(state.RaceStatus, state.RaceStatusPaused)
			c.Set(Confirmation, Unconfirmed)
			c.Set(Countdown, CountdownSteps)
			log.WithField("race-state", c.Get(state.RaceStatus)).
				WithField(Confirmation, c.Get(Confirmation)).
				Info("race started, ready for confirmation and countdown")
		}

		if s == Running {
			c.Set(state.RaceStatus, state.RaceStatusRunning)
		}

		if s == Stopped {
			c.Set(state.RaceStatus, state.RaceStatusStopped)
			c.Set(Mode, None)
			c.Set(StageNumber, c.Get(StageNumber).(uint8)+1)
			log.WithField("race-state", c.Get(state.RaceStatus)).
				WithField(Confirmation, c.Get(Confirmation)).
				Info("race stopped")
		}
	}
}

func (r *Rule) handleCountdown(c *state.Course, v *state.Value) {
	if s, ok := v.Get().(string); ok && v.Name() == Confirmation {
		if s == Confirmed {
			log.WithField("race-state", c.Get(state.RaceStatus)).
				Info("race confirmed")
			c.Set(Countdown, CountdownSteps)
			go r.countdown(c)
			log.WithField("race-state", c.Get(state.RaceStatus)).
				Info("race pre-start countdown started")
		}
	}

	if s, ok := v.Get().(uint8); ok && v.Name() == Countdown {
		if s <= 0 {
			c.Set(ActiveView, Running)
			log.WithField("race-state", c.Get(state.RaceStatus)).
				WithField("countdown", s).
				Info("race pre-start countdown complete")
		} else {
			log.WithField("race-state", c.Get(state.RaceStatus)).
				WithField("countdown", s).
				Info("race pre-start countdown")
		}
	}
}

//
// Handle car ready related events
//
func (r *Rule) handleReady(c *state.Car, v *state.Value) {
	if r.course.Get(ActiveView) == Started {
		if r.course.Get(Confirmation) == Unconfirmed {
			if v.Name() == state.ControllerBtnTrackCall && v.Get().(bool) {
				log.WithField("race-state", r.course.Get(state.RaceStatus)).
					WithField("car", c.Id()).
					Info("track call button pressed")
				c.Set(Ready, true)
			}
			if v.Name() == Ready && v.Get().(bool) {
				log.WithField("race-state", r.course.Get(state.RaceStatus)).
					WithField("car", c.Id()).
					Info("car entered ready state")
				var ready = uint8(0)
				for _, cr := range r.ready {
					if cr.Get(Ready).(bool) {
						ready++
					}
				}
				r.course.Set(CarsReady, ready)
			}
		}
	}
}

func (r *Rule) InitializeCarState(c *state.Car) {
	r.ready[c.Id()] = c
	r.course.Set(Cars, r.course.Get(Cars).(uint8)+1)

	c.Set(Ready, false)
	c.Subscribe(Ready, r)
	c.Subscribe(state.CarLap, r)
	c.Subscribe(state.ControllerBtnTrackCall, r)
}

func (r *Rule) InitializeCourseState(c *state.Course) {
	c.Set(Countdown, uint8(0))
	c.Set(Mode, None)
	c.Set(ActiveView, None)
	c.Set(Stage, r.stages[0])
	c.Set(StageNumber, uint8(0))
	c.Set(Cars, uint8(0))
	c.Set(CarsReady, uint8(0))
	c.Set(Confirmation, Unconfirmed)
	c.Subscribe(ActiveView, r)
	c.Subscribe(Countdown, r)
	c.Subscribe(Confirmation, r)
	c.Subscribe(CarsReady, r)
	c.Subscribe(StageNumber, r)
	c.Subscribe(state.RMSStatus, r)
	c.Subscribe(state.RaceTime, r)
	r.course = c
}

func (r *Rule) countdown(c *state.Course) {
	var countdown uint8
	for {
		select {
		case <-time.After(1 * time.Second):
			countdown = c.Get(Countdown).(uint8) - 1
			c.Set(Countdown, countdown)
		}
		if countdown <= 0 {
			return
		}
	}
}

func (s StageConfig) Metrics() []telemetry.Metric {
	return []telemetry.Metric{
		{Name: "race-stage-laps", Value: s.Laps},
		{Name: "race-stage-duration", Value: s.duration},
		{Name: "race-stage-type", Value: s.Type},
	}
}
*/
