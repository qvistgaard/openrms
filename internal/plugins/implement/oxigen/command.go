package oxigen

import (
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
)

const (
	CarForceLaneChangeLeft  = 0x80
	CarForceLaneChangeRight = 0x40
	CarForceLaneChangeNone  = 0x00
	CarForceLangeChangeAny  = CarForceLaneChangeLeft | CarForceLaneChangeRight
)

type Command struct {
	settings Settings
	state    byte
	car      *Car
}

type Car struct {
	id      uint8
	command byte
	value   byte
}

func newPitLaneSpeed(id uint8, speed uint8) *Car {
	return &Car{
		id:      id,
		command: 0x01,
		value:   speed,
	}
}

func newMaxSpeed(id uint8, speed state.MaxSpeed) *Car {
	return &Car{
		id:      id,
		command: 0x02,
		value:   uint8(speed),
	}
}

func newMinSpeed(id uint8, speed uint8, forceLC byte) *Car {
	return &Car{
		id:      id,
		command: 0x03,
		value:   (speed / 4) | forceLC,
	}
}

func newMaxBreaking(id uint8, maxBreaking uint8) *Car {
	return &Car{
		id:      id,
		command: 0x05,
		value:   maxBreaking,
	}
}

func newEmptyCommand(race state.CourseChanges, currentState byte, settings *Settings) *Command {
	log.Infof("currentState: %+v", currentState)
	c := &Command{
		state: currentState,
		settings: Settings{
			maxSpeed: settings.maxSpeed,
			pitLane: PitLane{
				lapCounting: settings.pitLane.lapCounting,
				lapTrigger:  settings.pitLane.lapTrigger,
			},
		},
	}
	for _, v := range race.Changes {
		log.Infof("Race change %+v", v)
		switch v.Name {
		case state.RaceStatus:
			bv := v.Value.(uint8)
			switch bv {
			case state.RaceStatusStopped:
				c.stop()
			case state.RaceStatusPaused:
				c.pause()
			case state.RaceStatusRunning:
				c.start()
			case state.RaceStatusFlaggedLCDisabled:
				c.flag(false)
			case state.RaceStatusFlaggedLCEnabled:
				c.flag(true)
			}
		case state.CourseMaxSpeed:
			bv := v.Value.(uint8)
			c.maxSpeed(bv)
		}
	}
	return c
}

func newSettings() *Settings {
	return &Settings{
		maxSpeed: 255,
		pitLane: PitLane{
			lapCounting: 0,
			lapTrigger:  0,
		},
	}
}

func (c *Command) carCommand(id uint8, s string, v interface{}) bool {
	switch s {
	case state.CarMaxSpeed:
		c.car = newMaxSpeed(id, v.(state.MaxSpeed))
	case state.CarMaxBreaking:
		c.car = newMaxBreaking(id, v.(uint8))
	case state.CarMinSpeed:
		c.car = newMinSpeed(id, v.(uint8), CarForceLaneChangeNone)
	case state.CarPitLaneSpeed:
		c.car = newPitLaneSpeed(id, v.(uint8))
	}
	if c.car != nil {
		return true
	} else {
		return false
	}

}

func (c *Command) maxSpeed(speed uint8) {
	c.settings.maxSpeed = speed
	log.WithField("max-speed", c.state).Debug("oxigen max car speed set.")

}

func (c *Command) start() {
	c.state = 0x03
	log.WithField("state", c.state).Debug("oxigen race state set to started.")

}

func (c *Command) pitLaneLapCount(enabled bool, entry bool) {
	if !enabled {
		c.settings.pitLane.lapCounting = 0x20
		c.settings.pitLane.lapTrigger = 0x00
		log.WithField("lap-counting", c.settings.pitLane.lapCounting).
			WithField("lap-trigger-on-entry", c.settings.pitLane.lapTrigger).
			Debug("oxigen pit lane lap counting disabled.")
	} else {
		c.settings.pitLane.lapCounting = 0x00
		if entry {
			c.settings.pitLane.lapTrigger = 0x00
		} else {
			c.settings.pitLane.lapTrigger = 0x40
		}
		log.WithField("lap-counting", c.settings.pitLane.lapCounting).
			WithField("lap-trigger-on-entry", c.settings.pitLane.lapTrigger).
			Debug("oxigen pit lane lap counting enabled.")
	}
}

func (c *Command) stop() {
	c.state = 0x01
	log.WithField("state", c.state).Debug("oxigen race state set to stopped.")
}

func (c *Command) pause() {
	c.state = 0x04
	log.WithField("state", c.state).Debug("oxigen race state set to paused.")
}

func (c *Command) flag(lc bool) {
	if lc {
		c.state = 0x05
		log.WithField("state", c.state).Debug("oxigen race state set to flagged with lane change enabled.")
	} else {
		c.state = 0x15
		log.WithField("state", c.state).Debug("oxigen race state set to flagged with lane change disabled.")
	}
}
