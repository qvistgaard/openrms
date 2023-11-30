package drivers

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"time"
)

// TODO go can require more then one interface as agument
type Car interface {
	SetMaxSpeed(percent uint8)
	SetPitLaneMaxSpeed(percent uint8)
	SetMaxBreaking(percent uint8)
	SetMinSpeed(percent uint8)
	Id() types.Id
	// Reset() bool
	// InPit() bool
	// Deslotted() bool
	// Controller() Controller
	// Lap() Lap
}

type Lap interface {
	Number() uint16
	Time() time.Duration
	Recorded() time.Duration
}

type Controller interface {
	BatteryWarning() bool
	Link() bool
	TrackCall() bool
	ArrowUp() bool
	ArrowDown() bool
	TriggerValue() float64
}

type Event interface {
	Car() Car
}

type PitLaneLapCounting int

const (
	LapCountingOnEntry PitLaneLapCounting = iota
	LapCountingOnExit
)

type PitLane interface {
	LapCounting(enabled bool, option PitLaneLapCounting)
}

type Track interface {
	MaxSpeed(percent uint8)
	PitLane() PitLane
}

type Race interface {
	Start()
	Flag()
	Pause()
	Stop()
}

type Driver interface {
	EventLoop() error
	EventChannel() <-chan Event

	Car(car types.Id) Car

	Track() Track
	Race() Race
	Init(ctx context.Context)
}
