package types

type CarPitState uint8

const (
	PitStateNotInPitLane CarPitState = iota
	PitStateEntered
	PitStateWaiting
	PitStateActive
	PitStateComplete
)
