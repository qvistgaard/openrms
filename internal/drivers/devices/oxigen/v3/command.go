package v3

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

func (c Command) pack(timer []byte, race *Race, track *Track) []byte {
	var cmd byte = 0x00
	var parameter byte = 0x00
	var controller byte = 0x00

	cmd = c.code
	parameter = c.value
	if c.id != nil {
		controller = *c.id
		cmd = 0x80 | cmd
	} else {
		controller = 0x00
		cmd = 0x00 | cmd
	}

	return []byte{
		race.status | track.pitLane.lapCounting | track.pitLane.lapCountingOption, //  0 o.race.status | o.track.pitLane.lapCounting | o.track.pitLane.lapCountingOption,
		track.maxSpeed, //  1 o.track.maxSpeed,
		controller,     //  2 link ID
		cmd,            //  3 Command value
		parameter,      //  4 command argument
		0x00,           //  5 unused
		0x00,           //  6 unused
		0x00,           //  7 unused
		timer[0],       //  8 Race timer
		timer[1],       //  9 Race timer
		timer[2],       // 10 Race timer
	}
}
