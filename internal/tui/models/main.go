package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/types"
	"golang.org/x/term"
	"os"
)

type ActiveView int

const (
	ViewLeaderboard ActiveView = iota
	ViewCarConfiguration
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
	Leaderboard      tea.Model
	CarConfiguration tea.Model
	raceControl      tea.Model
}

func (m Main) Init() tea.Cmd {
	// m.Leaderboard.Init()
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
			m.Bridge <- commands.StartRace{}
		case "p":
			m.Bridge <- commands.PauseRace{}
		case "r":
			m.Bridge <- commands.ResetRace{}
		}
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		if m.ActiveView == ViewCarConfiguration {
			m.CarConfiguration, cmd = m.CarConfiguration.Update(msg)
		}
	case types.RaceTelemetry:
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		//
	case tea.WindowSizeMsg:
		if m.ActiveView == ViewLeaderboard {
			m.Leaderboard, cmd = m.Leaderboard.Update(msg)
		}
		// m.Leaderboard, cmd = m.Leaderboard.Update(msg)

	case commands.OpenCarConfiguration:
		m.ActiveView = ViewCarConfiguration
		m.CarConfiguration, cmd = m.CarConfiguration.Update(msg)
	case commands.SaveCarConfiguration:
		m.ActiveView = ViewLeaderboard
		m.Bridge <- msg
	}

	return m, cmd
}

func (m Main) View() string {

	physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)
	if physicalWidth > 0 {
		docStyle = docStyle.Width(physicalWidth - 2)
	}
	if physicalWidth > 0 {
		docStyle = docStyle.Height(physicalHeight - 4)
	}

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		Header{}.View(),
		m.activeView(),

		//m.Leaderboard.View(),
	))
}

func (m Main) activeView() string {
	switch m.ActiveView {
	case ViewLeaderboard:
		return m.Leaderboard.View()
	case ViewCarConfiguration:
		return m.CarConfiguration.View()
	}
	return "No view"
}
