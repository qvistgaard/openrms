package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/elements"
)

var style1 = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("23")).
	Width(9)

var groupStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("250")).
	Width(7)

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
			style1.Width(h.width-76).Render("   ____                   ____  __  ________"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render("  / __ \\____  ___  ____  / __ \\/  |/  / ___/"),
			groupStyle.Render("Race: "),
			elements.Shortcut("R", "Start"),
			elements.Shortcut("F", "Flag"),
			elements.Shortcut("P", "Pause"),
			elements.Shortcut("S", "Stop"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render(" / / / / __ \\/ _ \\/ __ \\/ /_/ / /|_/ /\\__ \\ "),
			groupStyle.Render("Car:"),
			// helpStyle.Render("[ENTER] Details"),
			elements.Shortcut("C", "Configuration"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render("/ /_/ / /_/ /  __/ / / / _, _/ /  / /___/ / "),
			groupStyle.Render("Track:"),
			// helpStyle.Render("[ENTER] Details"),
			elements.Shortcut("T", "Configuration"),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render("\\____/ .___/\\___/_/ /_/_/ |_/_/  /_//____/  "),
		),
		lipgloss.JoinHorizontal(lipgloss.Right,
			style1.Render("    /_/                                     "),
		),
	)
}
