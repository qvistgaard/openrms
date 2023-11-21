package implement

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

type CarImplementer interface {
	MaxSpeed(percent types.Percent)
	PitLaneMaxSpeed(percent types.Percent)
	MaxBreaking(percent types.Percent)
	MinSpeed(percent types.Percent)
}

type PitLaneLapCounting int

const (
	LapCountingOnEntry PitLaneLapCounting = iota
	LapCountingOnExit
)

type PitLaneImplementer interface {
	LapCounting(enabled bool, option PitLaneLapCounting)
}

type TrackImplementer interface {
	MaxSpeed(percent types.Percent)
	PitLane() PitLaneImplementer
}

type RaceImplementer interface {
	Start()
	Flag()
	Pause()
	Stop()
}

type Implementer interface {
	EventLoop() error
	EventChannel() <-chan Event

	Car(car types.Id) CarImplementer

	Track() TrackImplementer
	Race() RaceImplementer
	Init(ctx context.Context, processor reactive.ValuePostProcessor)
}
