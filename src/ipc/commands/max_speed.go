package commands

import ipc "openrms/ipc"

type MaxSpeed struct {
	maxSpeed uint8
}

func NewMaxSpeed(driver uint8, maxSpeed uint8) *ipc.Command {
	c := &MaxSpeed{
		maxSpeed: maxSpeed,
	}
	return ipc.NewCommand(driver, c)
}

func (c *MaxSpeed) Value() []byte {
	return []byte{c.maxSpeed}
}
