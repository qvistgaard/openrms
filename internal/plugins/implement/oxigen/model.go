package oxigen

type Settings struct {
	maxSpeed byte
	pitLane  PitLane
}

type PitLane struct {
	lapCounting byte
	lapTrigger  byte
}
