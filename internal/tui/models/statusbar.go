package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/tui/messages"
	"time"
)

type StatusBar struct {
	width      int
	raceTime   time.Duration
	raceStatus race.Status
}

func (s StatusBar) Init() tea.Cmd {
	return nil
}

func (s StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.(tea.WindowSizeMsg).Width
	case messages.Update:
		s.raceTime = msg.(messages.Update).RaceDuration
		s.raceStatus = msg.(messages.Update).RaceStatus
	}
	return s, nil
}

func (s StatusBar) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("24")).
		PaddingLeft(1).
		PaddingRight(1).
		AlignHorizontal(lipgloss.Right).
		Width(s.width).
		Render("Race time: ", formatDuration(s.raceTime), "Status: ", formatRaceStatus(s.raceStatus))
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatRaceStatus(status race.Status) string {
	switch status {
	case race.Running:
		return "Running"
	case race.Stopped:
		return "Stopped"
	case race.Flagged:
		return "Flagged"
	case race.Paused:
		return "Paused"
	}
	return "N/A"
}
