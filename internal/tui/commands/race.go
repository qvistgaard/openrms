package commands

type StartRace struct {
	RaceTime string
	Laps     string
	// TODO implement laps and mapping to race plugin
}

type PauseRace struct {
}

type StopRace struct {
}

type FlagRace struct {
}

type ResumeRace struct {
}
