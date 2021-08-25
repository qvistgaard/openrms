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

func newPitLaneSpeed(id uint8, speed state.Speed) *Car {
	return &Car{
		id:      id,
		command: 0x81, // 0x81
		value:   state.PercentToUint8(float64(speed)),
	}
}

func newMaxSpeed(id uint8, speed state.Speed) *Car {
	return &Car{
		id:      id,
		command: 0x82,
		value:   state.PercentToUint8(float64(speed)),
	}
}

func newMinSpeed(id uint8, speed state.Speed, forceLC byte) *Car {
	return &Car{
		id:      id,
		command: 0x03,
		value:   (state.PercentToUint8(float64(speed)) / 4) | forceLC,
	}
}

func newMaxBreaking(id uint8, maxBreaking state.Breaking) *Car {
	return &Car{
		id:      id,
		command: 0x05,
		value:   state.PercentToUint8(float64(maxBreaking)),
	}
}

func newEmptyCommand(race state.CourseState, currentState byte, settings *Settings) *Command {
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

	lapChanges := false
	lapCounting := settings.pitLane.lapCounting != 0x20
	lapTrigger := settings.pitLane.lapTrigger != 0x00
	for _, v := range race.Changes {
		switch v.Name {
		case state.RaceStatus:
			if bv, ok := v.Value.(uint8); ok {
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
				log.WithField(state.RaceStatus, v).
					Debug("oxigen: requested new race-status")
			} else {
				log.WithField(state.RaceStatus, v).
					Warn("oxigen: discarded race-status, invalid value")
			}
		case state.CourseMaxSpeed:
			if bv, ok := v.Value.(state.Speed); ok {
				c.maxSpeed(bv)
				log.WithField(state.CourseMaxSpeed, v).
					Debug("oxigen: requested new course-max-speed")
			} else {
				log.WithField(state.CourseMaxSpeed, v).
					Warn("oxigen: discarded course-max-speed, invalid value")
			}
		case state.PitlaneLapCounting:
			if bv, ok := v.Value.(bool); ok {
				lapCounting = bv
				lapChanges = true
				log.WithField(state.PitlaneLapCounting, v).
					Debug("oxigen: requested new pitlane lap counting settings")
			} else {
				log.WithField(state.PitlaneLapCounting, v).
					Warn("oxigen: discarded pitlane lap counting, invalid value")
			}
		case state.PitlaneLapCountingOnEntry:
			if bv, ok := v.Value.(bool); ok {
				lapTrigger = bv
				lapChanges = true
				log.WithField(state.PitlaneLapCountingOnEntry, v).
					Debug("oxigen: requested new pitlane lap counting on entry settings")
			} else {
				log.WithField(state.PitlaneLapCountingOnEntry, v).
					Warn("oxigen: discarded pitlane lap counting on entry, invalid value")
			}

		}
	}
	if lapChanges {
		c.pitLaneLapCount(lapCounting, lapTrigger)
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
		if speed, ok := v.(state.Speed); ok {
			c.car = newMaxSpeed(id, speed)
			log.WithField("car", id).
				WithField("max-speed", v).
				Debugf("oxigen: new car max speed requested")
		} else {
			log.WithField("car", id).
				WithField("max-speed", v).
				Warn("oxigen: discarded car max speed command")
		}
	case state.CarMaxBreaking:
		if uintv, ok := v.(state.Breaking); ok {
			c.car = newMaxBreaking(id, uintv)
			log.WithField("car", id).
				WithField("max-breaking", v).
				Debugf("oxigen: new car max breaking requested")
		} else {
			log.WithField("car", id).
				WithField("max-breaking", v).
				Warn("oxigen: discarded car max breaking command")

		}
	case state.CarMinSpeed:
		if uintv, ok := v.(state.Speed); ok {
			c.car = newMinSpeed(id, uintv, CarForceLaneChangeNone)
			log.WithField("car", id).
				WithField("min-speed", v).
				WithField("force-lc", CarForceLaneChangeNone).
				Debugf("oxigen: new car min speed requested")
		} else {
			log.WithField("car", id).
				WithField("min-speed", v).
				WithField("force-lc", CarForceLaneChangeNone).
				Warn("oxigen: discarded min speed command")
		}
	case state.CarPitLaneSpeed:
		if speed, ok := v.(state.Speed); ok {
			c.car = newPitLaneSpeed(id, speed)
			log.WithField("car", id).
				WithField("max-pit-speed", v).
				Debugf("oxigen: new car max pit speed requested")
		} else {
			log.WithField("car", id).
				WithField("max-pit-speed", v).
				Warn("oxigen: discarded car max pit speed command")

		}
	}
	if c.car != nil {
		return true
	} else {
		return false
	}

}

func (c *Command) maxSpeed(speed state.Speed) {
	c.settings.maxSpeed = state.PercentToUint8(float64(speed))
}

func (c *Command) start() {
	c.state = 0x03
	log.WithField("state", c.state).Debug("oxigen: race state set to started.")
}

func (c *Command) pitLaneLapCount(enabled bool, entry bool) {
	if !enabled {
		c.settings.pitLane.lapCounting = 0x20
		c.settings.pitLane.lapTrigger = 0x00
		log.WithField("lap-counting", c.settings.pitLane.lapCounting).
			WithField("lap-trigger-on-entry", c.settings.pitLane.lapTrigger).
			Debug("oxigen: pit lane lap counting disabled.")
	} else {
		c.settings.pitLane.lapCounting = 0x00
		if entry {
			c.settings.pitLane.lapTrigger = 0x00
		} else {
			c.settings.pitLane.lapTrigger = 0x40
		}
		log.WithField("lap-counting", c.settings.pitLane.lapCounting).
			WithField("lap-trigger-on-entry", c.settings.pitLane.lapTrigger).
			Debug("oxigen: pit lane lap counting enabled.")
	}
}

func (c *Command) stop() {
	c.state = 0x01
	log.WithField("state", c.state).Debug("oxigen: race state set to stopped.")
}

func (c *Command) pause() {
	c.state = 0x04
	log.WithField("state", c.state).Debug("oxigen: race state set to paused.")
}

func (c *Command) flag(lc bool) {
	if lc {
		c.state = 0x05
		log.WithField("state", c.state).Debug("oxigen: race state set to flagged with lane change enabled.")
	} else {
		c.state = 0x15
		log.WithField("state", c.state).Debug("oxigen: race state set to flagged with lane change disabled.")
	}
}
