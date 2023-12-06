package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/tui/style"
)

type Confirmation struct {
	width      int
	height     int
	percentage float64
	progress   progress.Model
}

func InitialConfirmationModel() Confirmation {
	m := Confirmation{
		progress: progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
	}
	return m
}

func (c Confirmation) Init() tea.Cmd {
	//TODO drivers me
	panic("drivers me")
}

func (c Confirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height - 6
		c.progress.Width = 72
		return c, nil
	case confirmation.Status:
		c.percentage = 1 - (msg.RemainingTime.Seconds() / msg.TotalTime.Seconds())
		return c, nil
	}
	return c, nil
}

func (c Confirmation) View() string {
	return lipgloss.Place(c.width, c.height,
		lipgloss.Center, lipgloss.Center,
		style.DialogBox.Width(77).Height(0).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Waiting for confirmation"),
						c.progress.ViewAs(c.percentage),
						lipgloss.NewStyle().Width(72).MarginTop(1).MarginBottom(0).Align(lipgloss.Center).Render("10 / 10"),
					),
				),
			),
		),
	)
}
