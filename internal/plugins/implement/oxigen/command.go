package oxigen

import (
	"github.com/qvistgaard/openrms/internal/state"
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

func newMaxSpeed(id uint8, speed uint8) *Car {
	return &Car{
		id:      id,
		command: 0x02,
		value:   speed,
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

func newEmptyCommand(race map[string]state.StateInterface, currentState byte, settings *Settings) *Command {
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
	for k, v := range race {
		bv := v.Get().(uint8)
		switch k {
		case state.RaceStatus:
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
		case state.RaceMaxSpeed:
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

func (c *Command) carCommand(id uint8, s string, v state.StateInterface) bool {
	switch s {
	case state.CarMaxSpeed:
		c.car = newMaxSpeed(id, v.Get().(uint8))
	case state.CarMaxBreaking:
		c.car = newMaxBreaking(id, v.Get().(uint8))
	case state.CarMinSpeed:
		c.car = newMinSpeed(id, v.Get().(uint8), CarForceLaneChangeNone)
	case state.CarPitLaneSpeed:
		c.car = newPitLaneSpeed(id, v.Get().(uint8))
	}
	if c.car != nil {
		return true
	} else {
		return false
	}

}

func (c *Command) maxSpeed(speed uint8) {
	c.settings.maxSpeed = speed
}

func (c *Command) start() {
	c.state = 0x03
}

func (c *Command) pitLaneLapCount(enabled bool, entry bool) {
	if !enabled {
		c.settings.pitLane.lapCounting = 0x20
		c.settings.pitLane.lapTrigger = 0x00
	} else {
		c.settings.pitLane.lapCounting = 0x00
		if entry {
			c.settings.pitLane.lapTrigger = 0x00
		} else {
			c.settings.pitLane.lapTrigger = 0x40
		}
	}
}

func (c *Command) stop() {
	c.state = 0x01
}

func (c *Command) pause() {
	c.state = 0x04
}

func (c *Command) flag(lc bool) {
	if lc {
		c.state = 0x05
	} else {
		c.state = 0x15
	}
}
