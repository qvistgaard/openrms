package commands

type StartRace struct {
	RaceTime string
	Laps     string
}

type PauseRace struct {
}

type StopRace struct {
}

type FlagRace struct {
}

type ResumeRace struct {
}
