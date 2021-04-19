package state

type Rule interface {
	InitializeCarState(car *Car)
	InitializeRaceState(race *Race)
}
