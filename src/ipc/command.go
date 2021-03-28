package ipc

func NewCommand(driver uint8, command CommandType) *Command {
	c := new(Command)
	c.driver = driver
	c.command = command
	return c
}

func NewEmptyCommand() *Command {
	c := new(Command)
	c.driver = 0
	return c
}

type Command struct {
	driver  uint8
	command CommandType
}

func (c *Command) Driver() uint8 {
	return c.driver
}

func (c *Command) Value() []byte {
	return c.command.Value()
}

func (c *Command) CommandType() interface{} {
	return c.command
}

type CommandType interface {
	Value() []byte
}
