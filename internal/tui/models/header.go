package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style1 = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("23")).
	Width(10)

var groupStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("15")).
	Width(6)

var helpStyle = lipgloss.NewStyle().
	Width(20)

var headingStyle = lipgloss.NewStyle().
	Width(20).
	BorderForeground(lipgloss.Color("240")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true)

type Header struct {
}

func (h Header) Init() tea.Cmd {
	return nil
}

func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h Header) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render("OpenRMS"),
			groupStyle.Render("Race: "),
			helpStyle.Render("[S] Start"),
			helpStyle.Render("[P] Pause"),
			helpStyle.Render("[R] Reset"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render(""),
			groupStyle.Render("Car:"),
			// helpStyle.Render("[ENTER] Details"),
			helpStyle.Render("[C] Configuration"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render(""),
			groupStyle.Render(""),
			helpStyle.Render(""),
			helpStyle.Render(""),
			helpStyle.Render(""),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render(""),
			groupStyle.Render(""),
			helpStyle.Render(""),
			helpStyle.Render("Race Timer"),
			helpStyle.Render("00:00:00"),
		),
	)
}
