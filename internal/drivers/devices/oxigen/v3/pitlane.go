package v3

import (
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/rs/zerolog"
)

type PitLane struct {
	lapCounting       byte
	lapCountingOption byte
	logger            zerolog.Logger
}

func NewPitLane(logger zerolog.Logger) *PitLane {
	return &PitLane{
		logger:      logger,
		lapCounting: pitLaneLapCountingEnabledByte,
	}
}

const (
	pitLaneLapCountingEnabledByte  = 0x00
	pitLaneLapCountingDisabledByte = 0x20
	pitLaneLapCountingOnEntryByte  = 0x00
	pitLaneLapCountingOnExitByte   = 0x40
)

func (p *PitLane) LapCounting(enabled bool, option drivers.PitLaneLapCounting) {
	if !enabled {
		p.lapCounting = pitLaneLapCountingDisabledByte
		p.lapCountingOption = pitLaneLapCountingOnEntryByte
	} else {
		p.lapCounting = pitLaneLapCountingDisabledByte
		switch option {
		case drivers.LapCountingOnEntry:
			p.lapCountingOption = pitLaneLapCountingOnEntryByte
		case drivers.LapCountingOnExit:
			p.lapCountingOption = pitLaneLapCountingOnExitByte
		}
	}
}
