package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/state/race"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/messages"
)

type ActiveView int

const (
	ViewLeaderboard ActiveView = iota
	ViewCarConfiguration
	ViewRaceConfiguration
	ViewTrackConfiguration
)

type Main struct {
	Bridge             chan<- tea.Msg
	ActiveView         ActiveView
	Header             tea.Model
	StatusBar          tea.Model
	Leaderboard        tea.Model
	CarConfiguration   tea.Model
	RaceControl        tea.Model
	width              int
	height             int
	raceStatus         race.RaceStatus
	TrackConfiguration tea.Model
}

func (m Main) Init() tea.Cmd {
	m.Leaderboard.Init()
	m.StatusBar.Init()
	m.CarConfiguration.Init()
	return nil
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case ActiveView:
		m.ActiveView = msg.(ActiveView)
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			return m, tea.Quit
		case "r":
			if m.raceStatus == race.RaceStopped {
				m.ActiveView = ViewRaceConfiguration
				return m, nil
			} else {
				m.Bridge <- commands.ResumeRace{}
			}
		case "p":
			m.Bridge <- commands.PauseRace{}
		case "f":
			m.Bridge <- commands.FlagRace{}
		case "s":
			m.Bridge <- commands.StopRace{}
		case "t":
			return m, func() tea.Msg {
				return commands.OpenTrackConfiguration{}
			}
		}
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		if m.ActiveView == ViewCarConfiguration {
			m.CarConfiguration, cmd = m.CarConfiguration.Update(msg)
		}
		if m.ActiveView == ViewRaceConfiguration {
			m.RaceControl, cmd = m.RaceControl.Update(msg)
		}
		if m.ActiveView == ViewTrackConfiguration {
			m.TrackConfiguration, cmd = m.TrackConfiguration.Update(msg)
		}
	case messages.Update:
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		m.raceStatus = msg.(messages.Update).RaceStatus
		m.StatusBar, _ = m.StatusBar.Update(msg)
		m.TrackConfiguration, _ = m.TrackConfiguration.Update(msg)

	case tea.WindowSizeMsg:
		width := msg.(tea.WindowSizeMsg).Width
		height := msg.(tea.WindowSizeMsg).Height
		updatedMsg := tea.WindowSizeMsg{
			Width:  width - 4,
			Height: height - 2,
		}
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(updatedMsg)
		}
		m.height = height
		m.width = width
		m.StatusBar, _ = m.StatusBar.Update(updatedMsg)
		m.Header, _ = m.Header.Update(updatedMsg)
		m.CarConfiguration, _ = m.CarConfiguration.Update(updatedMsg)
		m.RaceControl, _ = m.RaceControl.Update(updatedMsg)
		m.TrackConfiguration, _ = m.TrackConfiguration.Update(updatedMsg)

	case commands.OpenCarConfiguration:
		m.ActiveView = ViewCarConfiguration
		m.CarConfiguration, cmd = m.CarConfiguration.Update(msg)
	case commands.OpenTrackConfiguration:
		m.ActiveView = ViewTrackConfiguration
		m.TrackConfiguration, cmd = m.TrackConfiguration.Update(msg)
	case commands.SaveCarConfiguration:
		m.ActiveView = ViewLeaderboard
		m.Bridge <- msg
	case commands.SaveTrackConfiguration:
		m.ActiveView = ViewLeaderboard
		m.Bridge <- msg
	case commands.StartRace:
		m.ActiveView = ViewLeaderboard
		m.RaceControl, cmd = m.RaceControl.Update(msg)
		m.Bridge <- msg
	}

	return m, cmd
}

func (m Main) View() string {
	docStyle := lipgloss.NewStyle().Padding(0, 2, 1, 2)
	if m.width > 0 {
		docStyle = docStyle.Width(m.width)
	}
	if m.height > 0 {
		docStyle = docStyle.Height(m.height)
	}

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		m.Header.View(),
		m.activeView(),
		m.StatusBar.View(),
	))
}

func (m Main) activeView() string {
	switch m.ActiveView {
	case ViewLeaderboard:
		return m.Leaderboard.View()
	case ViewCarConfiguration:
		return m.CarConfiguration.View()
	case ViewRaceConfiguration:
		return m.RaceControl.View()
	case ViewTrackConfiguration:
		return m.TrackConfiguration.View()
	}
	return "No view"
}
