package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madflojo/tasks"
	"github.com/qvistgaard/openrms/internal/plugins/postprocessors/leaderboard"
	"github.com/qvistgaard/openrms/internal/repostitory/car"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

func Run(b *Bridge) {

	/*	c.Scheduler.Add(&tasks.Task{
			Interval: 1 * time.Second,
			TaskFunc: func() error {
				board := b.Leaderboard
				if board != nil {
					p.Send(board)
				}
				return nil
			},
		})
	*/
}

type Bridge struct {
	Leaderboard   *leaderboard.Leaderboard
	Scheduler     *tasks.Scheduler
	RaceTelemetry types.RaceTelemetry
	Program       *tea.Program
	UI            *UI
	messages      <-chan tea.Msg
	Cars          car.Repository
}

func CreateBridge(leaderboard *leaderboard.Leaderboard, scheduler *tasks.Scheduler, cars car.Repository) *Bridge {
	bridgeChannel := make(chan tea.Msg)

	return &Bridge{
		messages:    bridgeChannel,
		Leaderboard: leaderboard,
		Scheduler:   scheduler,
		Cars:        cars,
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
				bridge.UI.Send(bridge.RaceTelemetry)
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
				car.MaxSpeed().Set(maxSpeed)
			}
			log.Info(msg)
		}
	}
}
