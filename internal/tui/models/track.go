package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/elements"
	"github.com/qvistgaard/openrms/internal/tui/style"
)

type TrackControl struct {
	focusIndex    int
	focusIndexMax int
	inputs        []*textinput.Model
	maxSpeed      textinput.Model
	width         int
	height        int
}

func InitialTrackControlModel() TrackControl {
	m := TrackControl{}

	m.maxSpeed = textinput.New()
	m.maxSpeed.PromptStyle = style.Form.PromptStyle.Focused.Copy().Width(18)
	m.maxSpeed.Placeholder = "100"
	m.maxSpeed.Prompt = "Max speed:"
	m.maxSpeed.Focus()
	m.maxSpeed.TextStyle = style.Form.TextStyle.Focused.Copy()
	m.maxSpeed.CharLimit = 3

	m.inputs = []*textinput.Model{&m.maxSpeed}
	m.focusIndexMax = len(m.inputs) + 1
	return m
}

func (r TrackControl) Init() tea.Cmd {
	//TODO drivers me
	panic("drivers me")
}

func (r TrackControl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.inputs = []*textinput.Model{&r.maxSpeed}

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
					return commands.SaveTrackConfiguration{
						MaxSpeed: r.maxSpeed.Value(),
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
				(*input).PromptStyle = style.Form.PromptStyle.Blurred.Copy().Width(18)
				(*input).TextStyle = style.Form.TextStyle.Blurred.Copy()
			}

			if r.focusIndex >= 0 && r.focusIndex < len(r.inputs) {
				(r.inputs[r.focusIndex]).PromptStyle = style.Form.PromptStyle.Focused.Copy().Width(18)
				(r.inputs[r.focusIndex]).TextStyle = style.Form.TextStyle.Focused.Copy()
				cmds[r.focusIndex] = (r.inputs[r.focusIndex]).Focus()
			}
			return r, tea.Batch(cmds...)

		}
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height - 6
		return r, nil
	case commands.OpenTrackConfiguration:
		r.focusIndex = 0
		r.maxSpeed.Focus()
		r.maxSpeed.SetValue(msg.MaxSpeed)

		return r, nil
	}

	var cmd tea.Cmd
	if r.focusIndex >= 0 && r.focusIndex < len(r.inputs) {
		*r.inputs[r.focusIndex], cmd = (*r.inputs[r.focusIndex]).Update(msg)
	}
	return r, cmd
}

func (r TrackControl) View() string {
	ok := elements.Button("Save", r.focusIndex == len(r.inputs))
	cancel := elements.Button("Cancel", r.focusIndex == r.focusIndexMax)

	return lipgloss.Place(r.width, r.height,
		lipgloss.Center, lipgloss.Center,
		style.DialogBox.Width(77).Height(3).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Track configuration"),
						r.maxSpeed.View(),
					),
				),
				lipgloss.PlaceHorizontal(78, lipgloss.Center,
					lipgloss.JoinHorizontal(lipgloss.Center, ok, cancel),
				),
			),
		),
	)
}
