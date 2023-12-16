package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
	"github.com/qvistgaard/openrms/internal/plugins/flags"
)

type SimpleConfirmation struct {
	percentage float64
	progress   progress.Model
	flag       flags.Flag
}

func InitialSimpleConfirmationModel() SimpleConfirmation {
	m := SimpleConfirmation{
		progress: progress.New(progress.WithSolidFill("24"), progress.WithoutPercentage()),
	}
	m.progress.Width = 72
	return m
}

func (c SimpleConfirmation) Init() tea.Cmd {
	//TODO drivers me
	panic("drivers me")
}

func (c SimpleConfirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case confirmation.Status:
		c.progress.Width = 72
		c.percentage = 1 - (msg.RemainingTime.Seconds() / msg.TotalTime.Seconds())
		if c.percentage >= 1 {
			c.percentage = -1
		}
		return c, nil
	case flags.Flag:
		c.flag = msg

	}
	return c, nil
}

func (c SimpleConfirmation) View() string {
	if c.percentage >= 0 {
		return c.progress.ViewAs(c.percentage)
	}
	switch c.flag {
	case flags.Yellow:
		return lipgloss.NewStyle().
			Width(72).Background(lipgloss.Color("220")).
			Foreground(lipgloss.Color("16")).
			AlignHorizontal(lipgloss.Center).
			Bold(true).
			Render("CAUTION")
	case flags.Green:
		return lipgloss.NewStyle().Width(72).Background(lipgloss.Color("70")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Render("ALL CLEAR")
	default:
		return lipgloss.NewStyle().Width(72).Background(lipgloss.Color("232")).Render("")
	}
}
