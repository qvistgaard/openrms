package implement

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
)

type CarImplementer interface {
	MaxSpeed(percent uint8)
	PitLaneMaxSpeed(percent uint8)
	MaxBreaking(percent uint8)
	MinSpeed(percent uint8)
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
	MaxSpeed(percent uint8)
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
	Init(ctx context.Context)
}
