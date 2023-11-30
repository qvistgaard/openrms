package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madflojo/tasks"
	"github.com/qvistgaard/openrms/internal/plugins/leaderboard"
	race2 "github.com/qvistgaard/openrms/internal/plugins/race"
	"github.com/qvistgaard/openrms/internal/state/car/repository"
	"github.com/qvistgaard/openrms/internal/state/observable"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/state/track"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/messages"
	"github.com/qvistgaard/openrms/internal/types"
	"strconv"
	"time"
)

type Bridge struct {
	Leaderboard   *leaderboard.Plugin
	Scheduler     *tasks.Scheduler
	RaceTelemetry types.RaceTelemetry
	Program       *tea.Program
	Race          *race.Race
	Track         *track.Track
	UI            *UI
	messages      <-chan tea.Msg
	Cars          repository.Repository
	duration      time.Duration
	status        race.RaceStatus
	racePlugin    *race2.Plugin
	trackMaxSpeed uint8
}

func CreateBridge(leaderboard *leaderboard.Plugin, plugin *race2.Plugin, scheduler *tasks.Scheduler, track *track.Track, cars repository.Repository, race *race.Race) *Bridge {
	bridgeChannel := make(chan tea.Msg)

	return &Bridge{
		messages:    bridgeChannel,
		Leaderboard: leaderboard,
		racePlugin:  plugin,
		duration:    time.Second * 0,
		Scheduler:   scheduler,
		Cars:        cars,
		Track:       track,
		Race:        race,
		UI:          Create(bridgeChannel),
	}
}

func (bridge *Bridge) Run() {
	go bridge.messageHandler()

	bridge.Race.Duration().RegisterObserver(func(duration time.Duration, annotations observable.Annotations) {
		bridge.duration = duration
	})
	bridge.Race.Status().RegisterObserver(func(status race.RaceStatus, annotations observable.Annotations) {
		bridge.status = status
	})
	bridge.Leaderboard.RegisterObserver(func(telemetry types.RaceTelemetry, annotations observable.Annotations) {
		bridge.RaceTelemetry = telemetry
	})
	bridge.Track.MaxSpeed().RegisterObserver(func(maxSpeed uint8, annotations observable.Annotations) {
		bridge.trackMaxSpeed = maxSpeed
	})

	bridge.Scheduler.Add(&tasks.Task{
		Interval: 1 * time.Second,
		TaskFunc: func() error {
			if bridge.RaceTelemetry != nil && bridge.UI != nil {
				bridge.UI.Send(messages.Update{
					RaceTelemetry: bridge.RaceTelemetry,
					RaceStatus:    bridge.status,
					RaceDuration:  bridge.duration,
					TrackMaxSpeed: bridge.trackMaxSpeed,
				})
			}
			return nil
		},
	})
}

func (bridge *Bridge) messageHandler() {
	for {
		select {
		case msg := <-bridge.messages:
			switch msg := msg.(type) {
			case commands.SaveCarConfiguration:
				bridge.saveCarConfiguration(msg)
			case commands.SaveTrackConfiguration:
				bridge.saveTrackConfiguration(msg)
			case commands.StartRace:
				bridge.Race.Start()
			case commands.ResumeRace:
				bridge.Race.Start()
			case commands.PauseRace:
				bridge.Race.Pause()
			case commands.StopRace:
				bridge.Race.Stop()
			case commands.FlagRace:
				bridge.Race.Flag()
			}
		}
	}
}

func (bridge *Bridge) saveTrackConfiguration(msg commands.SaveTrackConfiguration) {
	maxSpeed, _ := strconv.ParseUint(msg.MaxSpeed, 10, 8)
	bridge.Track.MaxSpeed().Set(uint8(maxSpeed))
}

func (bridge *Bridge) saveCarConfiguration(msg commands.SaveCarConfiguration) {
	fromString, _ := types.IdFromString(msg.CarId)
	car, _, _ := bridge.Cars.Get(fromString)
	maxSpeed, _ := strconv.ParseUint(msg.MaxSpeed, 10, 8)
	maxPitSpeed, _ := strconv.ParseUint(msg.MaxPitSpeed, 10, 8)
	minSpeed, _ := strconv.ParseUint(msg.MinSpeed, 10, 8)
	name := msg.DriverName

	car.MaxSpeed().Set(uint8(maxSpeed))
	car.PitLaneMaxSpeed().Set(uint8(maxPitSpeed))
	car.MinSpeed().Set(uint8(minSpeed))
	car.Drivers().Set(types.Drivers{
		{Name: name},
	})
}
