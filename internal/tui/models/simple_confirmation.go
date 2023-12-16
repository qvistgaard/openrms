package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qvistgaard/openrms/internal/plugins/confirmation"
)

type SimpleConfirmation struct {
	percentage float64
	progress   progress.Model
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
			c.percentage = 0
		}
		return c, nil
	}
	return c, nil
}

func (c SimpleConfirmation) View() string {
	return c.progress.ViewAs(c.percentage)
}
