package oxigen

import (
	"github.com/qvistgaard/openrms/internal/implement"
)

type PitLane struct {
	lapCounting       byte
	lapCountingOption byte
}

func NewPitLane() *PitLane {
	return &PitLane{
		lapCounting: pitLaneLapCountingEnabledByte,
	}
}

const (
	pitLaneLapCountingEnabledByte  = 0x00
	pitLaneLapCountingDisabledByte = 0x20
	pitLaneLapCountingOnEntryByte  = 0x00
	pitLaneLapCountingOnExitByte   = 0x40
)

func (p *PitLane) LapCounting(enabled bool, option implement.PitLaneLapCounting) {
	if !enabled {
		p.lapCounting = pitLaneLapCountingDisabledByte
		p.lapCountingOption = pitLaneLapCountingOnEntryByte
	} else {
		p.lapCounting = pitLaneLapCountingDisabledByte
		switch option {
		case implement.LapCountingOnEntry:
			p.lapCountingOption = pitLaneLapCountingOnEntryByte
		case implement.LapCountingOnExit:
			p.lapCountingOption = pitLaneLapCountingOnExitByte
		}
	}
}
