package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madflojo/tasks"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/leaderboard"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/messages"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/reactivex/rxgo/v2"
	"time"
)

type Bridge struct {
	Leaderboard   *leaderboard.Leaderboard
	Scheduler     *tasks.Scheduler
	RaceTelemetry types.RaceTelemetry
	Program       *tea.Program
	Race          *race.Race
	UI            *UI
	messages      <-chan tea.Msg
	Cars          car.Repository
}

func CreateBridge(leaderboard *leaderboard.Leaderboard, scheduler *tasks.Scheduler, cars car.Repository, race *race.Race) *Bridge {
	bridgeChannel := make(chan tea.Msg)

	return &Bridge{
		messages:    bridgeChannel,
		Leaderboard: leaderboard,
		Scheduler:   scheduler,
		Cars:        cars,
		Race:        race,
		UI:          Create(bridgeChannel),
	}
}

func (bridge *Bridge) Run() {
	go bridge.messageHandler()

	bridge.Leaderboard.RegisterObserver(func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			bridge.RaceTelemetry = i.(types.RaceTelemetry)
		})
	})
	bridge.Scheduler.Add(&tasks.Task{
		Interval: 1 * time.Second,
		TaskFunc: func() error {
			if bridge.RaceTelemetry != nil && bridge.UI != nil {
				bridge.UI.Send(messages.Update{
					RaceTelemetry: bridge.RaceTelemetry,
					RaceStatus:    bridge.Race.CurrentState(),
					RaceDuration:  bridge.Race.Duration(),
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
				fromString, _ := types.IdFromString(msg.CarId)
				car, _, _ := bridge.Cars.Get(fromString, context.TODO())
				maxSpeed, _ := types.PercentFromString(msg.MaxSpeed)
				maxPitSpeed, _ := types.PercentFromString(msg.MaxPitSpeed)
				minSpeed, _ := types.PercentFromString(msg.MinSpeed)
				name := msg.DriverName

				car.MaxSpeed().Set(maxSpeed)
				car.PitLaneMaxSpeed().Set(maxPitSpeed)
				car.MinSpeed().Set(minSpeed)
				car.Drivers().Set(types.Drivers{
					{Name: name},
				})
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
			// log.Info(fmt.Sprintf("%+v\n", msg))
		}
	}
}
