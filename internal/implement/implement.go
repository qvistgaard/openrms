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

type RaceStatus int

const (
	RaceStopped RaceStatus = iota
	RacePaused
	RaceRunning
	RaceFlagged
)

type PitLaneImplementer interface {
	LapCounting(enabled bool, option PitLaneLapCounting)
}

type TrackImplementer interface {
	MaxSpeed(percent types.Percent)
	PitLane() PitLaneImplementer
}

type RaceImplementer interface {
	Status(status RaceStatus)
}

type Implementer interface {
	EventLoop() error
	EventChannel() <-chan Event

	Car(car types.Id) CarImplementer

	Track() TrackImplementer
	Race() RaceImplementer
	Init(ctx context.Context, processor reactive.ValuePostProcessor)

	/*	CarMaxSpeed(car uint, percent state.Percent)
		CarMaxBreaking(car uint, percent state.Percent)
		CarMinSpeed(car uint, percent state.Percent)
		CarPitLaneSpeed(car uint, percent state.Percent)
	*/

	/*	CourseMaxSpeed(percent state.Percent)*/

	/*	SendCarState(c interface{}) error
		SendRaceState(r interface{}) error*/

	// Resend relevant car state to implement.
	//
	// this method is executed if for example the controller
	// looses link with the dongle. But also for each car if
	// race status changes.
	// ResendCarState(c interface{})
}
