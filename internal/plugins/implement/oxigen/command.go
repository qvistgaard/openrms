package oxigen

type Command struct {
	id    *uint8
	code  byte
	value byte
}

func newCommand(id *uint8, code byte, value byte) Command {
	return Command{id: id, code: code, value: value}
}

func newEmptyCommand() Command {
	return Command{
		id:    nil,
		code:  0x00,
		value: 0x00,
	}
}
