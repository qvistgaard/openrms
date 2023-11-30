package oxigen

import "github.com/qvistgaard/openrms/internal/drivers"

func Controller(data []byte) drivers.Controller {
	return controller{data}
}

type controller struct {
	data []byte
}

func (c controller) BatteryWarning() bool {
	return 0x04&c.data[0] == 0x04
}

func (c controller) Link() bool {
	return 0x02&c.data[0] == 0x02
}

func (c controller) TrackCall() bool {
	return 0x08&c.data[0] == 0x08
}

func (c controller) ArrowUp() bool {
	return 0x20&c.data[0] == 0x20
}

func (c controller) ArrowDown() bool {
	return 0x40&c.data[0] == 0x40
}

func (c controller) TriggerValue() float64 {
	return float64(0x7F & c.data[7])
}
