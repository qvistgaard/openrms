package models

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/elements"
	"github.com/qvistgaard/openrms/internal/tui/style"
	"strconv"
)

type RaceControl struct {
	focusIndex    int
	focusIndexMax int
	inputs        []*textinput.Model
	raceTime      textinput.Model
}

func InitialRaceControlModel() RaceControl {
	m := RaceControl{}

	m.raceTime = textinput.New()
	m.raceTime.Focus()
	m.raceTime.PromptStyle = focusedStyle.Copy().Width(18)
	m.raceTime.Prompt = "Race timer:"
	m.raceTime.TextStyle = focusedStyle
	m.raceTime.CharLimit = 64
	m.raceTime.Placeholder = "10m"

	m.raceTime.Cursor.SetMode(cursor.CursorBlink)

	m.inputs = []*textinput.Model{&m.raceTime}
	m.focusIndexMax = len(m.inputs) + 1
	return m
}

func (r RaceControl) Init() tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func (r RaceControl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.inputs = []*textinput.Model{&r.raceTime}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return r, func() tea.Msg {
				return ViewLeaderboard
			}

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && r.focusIndex == r.focusIndexMax {
				return r, func() tea.Msg {
					return ViewLeaderboard
				}
			}
			if s == "enter" && r.focusIndex == len(r.inputs) {
				return r, func() tea.Msg {
					return commands.StartRace{
						RaceTime: r.raceTime.Value(),
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				r.focusIndex--
			} else {
				r.focusIndex++
			}

			if r.focusIndex > r.focusIndexMax {
				r.focusIndex = 0
			} else if r.focusIndex < 0 {
				r.focusIndex = r.focusIndexMax
			}

			cmds := make([]tea.Cmd, len(r.inputs))

			for _, input := range r.inputs {
				(*input).Blur()
				(*input).PromptStyle = blurredStyle.Copy().Width(18)
			}

			if r.focusIndex >= 0 && r.focusIndex < len(r.inputs) {
				(r.inputs[r.focusIndex]).PromptStyle = focusedStyle.Copy().Width(18)
				cmds[r.focusIndex] = (r.inputs[r.focusIndex]).Focus()
			}
			return r, tea.Batch(cmds...)
		}
	}

	var cmd tea.Cmd
	if r.focusIndex >= 0 && r.focusIndex < len(r.inputs) {
		*r.inputs[r.focusIndex], cmd = (*r.inputs[r.focusIndex]).Update(msg)
	}
	return r, cmd
}

func (r RaceControl) View() string {
	ok := elements.Button("Start", r.focusIndex == len(r.inputs))
	cancel := elements.Button("Cancel", r.focusIndex == r.focusIndexMax)

	return lipgloss.PlaceHorizontal(200,
		lipgloss.Center,
		style.DialogBox.Width(77).Height(3).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Race configuration"),
						r.raceTime.View(),
					),
				),
				lipgloss.PlaceHorizontal(78, lipgloss.Center,
					lipgloss.JoinHorizontal(lipgloss.Center, ok, cancel, strconv.Itoa(r.focusIndex)),
				),
			),
		),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}
