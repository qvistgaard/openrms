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
)

var (

	// General.

	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	// Tabs.

	// Dialog.

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

type Main struct {
	Bridge           chan<- tea.Msg
	ActiveView       ActiveView
	Header           tea.Model
	StatusBar        tea.Model
	Leaderboard      tea.Model
	CarConfiguration tea.Model
	RaceControl      tea.Model
	width            int
	height           int
	raceStatus       race.RaceStatus
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
		case "s":
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
		case "e":
			m.Bridge <- commands.StopRace{}
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
	case messages.Update:
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		m.raceStatus = msg.(messages.Update).RaceStatus
		m.StatusBar, _ = m.StatusBar.Update(msg)

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

	case commands.OpenCarConfiguration:
		m.ActiveView = ViewCarConfiguration
		m.CarConfiguration, cmd = m.CarConfiguration.Update(msg)
	case commands.SaveCarConfiguration:
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

	// physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)
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
	}
	return "No view"
}
