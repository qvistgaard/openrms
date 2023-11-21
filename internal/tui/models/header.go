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
	width int
}

func (h Header) Init() tea.Cmd {
	return nil
}

func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.(tea.WindowSizeMsg).Width
	}
	return h, nil
}

func (h Header) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Width(h.width-76).Render("OpenRMS"),
			groupStyle.Render("Race: "),
			helpStyle.Render("[S] Start"),
			helpStyle.Render("[F] Flag"),
			helpStyle.Render("[P] Pause"),
			helpStyle.Render("[E] Stop"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render(""),
			groupStyle.Render("Car:"),
			// helpStyle.Render("[ENTER] Details"),
			helpStyle.Render("[C] Configuration"),
		),
	)
}
