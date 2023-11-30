package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qvistgaard/openrms/internal/tui/models"
	"os"
)

type UI struct {
	program *tea.Program
}

func Create(bridge chan<- tea.Msg) *UI {
	return &UI{
		program: tea.NewProgram(models.Main{
			Bridge:             bridge,
			Header:             models.Header{},
			StatusBar:          models.StatusBar{},
			ActiveView:         models.ViewLeaderboard,
			Leaderboard:        models.InitialLeaderboardModel(),
			CarConfiguration:   models.InitialCarConfigurationModel(),
			RaceControl:        models.InitialRaceControlModel(),
			TrackConfiguration: models.InitialTrackControlModel(),
		}, tea.WithAltScreen()),
	}
}

func (ui *UI) Run() {
	if _, err := ui.program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (ui *UI) Send(msg tea.Msg) {
	ui.program.Send(msg)
}
