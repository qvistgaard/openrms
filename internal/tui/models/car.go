package models

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/elements"
	"github.com/qvistgaard/openrms/internal/tui/style"
)

type CarDetails struct {
}

type CarConfiguration struct {
	car string

	focusIndex    int
	focusIndexMax int
	inputs        []*textinput.Model
	onTrackSpeed  textinput.Model
	inPitSpeed    textinput.Model
	teamName      textinput.Model
	minSpeed      textinput.Model
	width         int
	height        int
}

func InitialCarConfigurationModel() CarConfiguration {
	m := CarConfiguration{}

	m.teamName = textinput.New()
	m.teamName.Focus()
	m.teamName.PromptStyle = style.Form.PromptStyle.Focused.Copy().Width(18)
	m.teamName.Prompt = "Team name:"
	m.teamName.TextStyle = style.Form.TextStyle.Focused.Copy()
	m.teamName.CharLimit = 64
	m.teamName.Cursor.SetMode(cursor.CursorBlink)

	m.onTrackSpeed = textinput.New()
	m.onTrackSpeed.PromptStyle = style.Form.PromptStyle.Blurred.Copy().Width(18)
	m.onTrackSpeed.Placeholder = "100"
	m.onTrackSpeed.Prompt = "On track:"
	m.onTrackSpeed.TextStyle = style.Form.TextStyle.Blurred.Copy()
	m.onTrackSpeed.CharLimit = 3

	m.inPitSpeed = textinput.New()
	m.inPitSpeed.Placeholder = "75"
	m.inPitSpeed.PromptStyle = style.Form.PromptStyle.Blurred.Copy().Width(18)
	m.inPitSpeed.Prompt = "In pit lane:"
	m.inPitSpeed.TextStyle = style.Form.TextStyle.Blurred.Copy()
	m.inPitSpeed.CharLimit = 3

	m.minSpeed = textinput.New()
	m.minSpeed.Placeholder = "75"
	m.minSpeed.PromptStyle = style.Form.PromptStyle.Blurred.Copy().Width(18)
	m.minSpeed.Prompt = "Min speed:"
	m.minSpeed.TextStyle = style.Form.TextStyle.Blurred.Copy()
	m.minSpeed.CharLimit = 3

	m.inputs = []*textinput.Model{&m.teamName, &m.onTrackSpeed, &m.inPitSpeed, &m.minSpeed}
	m.focusIndexMax = len(m.inputs) + 1
	return m
}

func (c CarConfiguration) Init() tea.Cmd {
	return textinput.Blink
}

func (c CarConfiguration) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	c.inputs = []*textinput.Model{&c.teamName, &c.onTrackSpeed, &c.inPitSpeed, &c.minSpeed}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return c, func() tea.Msg {
				return ViewLeaderboard
			}

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && c.focusIndex == c.focusIndexMax {
				return c, func() tea.Msg {
					return ViewLeaderboard
				}
			}
			if s == "enter" && c.focusIndex == len(c.inputs) {
				return c, func() tea.Msg {
					return commands.SaveCarConfiguration{
						CarId:       c.car,
						MaxSpeed:    c.onTrackSpeed.Value(),
						MinSpeed:    c.minSpeed.Value(),
						MaxPitSpeed: c.inPitSpeed.Value(),
						DriverName:  c.teamName.Value(),
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			if c.focusIndex > c.focusIndexMax {
				c.focusIndex = 0
			} else if c.focusIndex < 0 {
				c.focusIndex = c.focusIndexMax
			}

			cmds := make([]tea.Cmd, len(c.inputs))

			for _, input := range c.inputs {
				(*input).Blur()
				(*input).PromptStyle = style.Form.PromptStyle.Blurred.Copy().Width(18)
				(*input).TextStyle = style.Form.TextStyle.Blurred.Copy()
			}

			if c.focusIndex >= 0 && c.focusIndex < len(c.inputs) {
				(c.inputs[c.focusIndex]).PromptStyle = style.Form.PromptStyle.Focused.Copy().Width(18)
				(c.inputs[c.focusIndex]).TextStyle = style.Form.TextStyle.Focused.Copy()
				cmds[c.focusIndex] = (c.inputs[c.focusIndex]).Focus()
			}
			return c, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height - 6
		return c, nil

	case commands.OpenCarConfiguration:
		c.focusIndex = 0
		c.car = msg.CarId
		c.teamName.SetValue(msg.DriverName)
		c.inPitSpeed.SetValue(msg.MaxPitSpeed)
		c.onTrackSpeed.SetValue(msg.MaxSpeed)
		c.minSpeed.SetValue(msg.MinSpeed)
		c.teamName.Focus()
		return c, nil
	}

	var cmd tea.Cmd
	if c.focusIndex >= 0 && c.focusIndex < len(c.inputs) {
		*c.inputs[c.focusIndex], cmd = (*c.inputs[c.focusIndex]).Update(msg)
	}
	return c, cmd
}

func (c CarConfiguration) View() string {
	ok := elements.Button("Save", c.focusIndex == len(c.inputs))
	cancel := elements.Button("Cancel", c.focusIndex == c.focusIndexMax)

	return lipgloss.Place(c.width, c.height,
		lipgloss.Center, lipgloss.Center,
		style.DialogBox.Width(77).Height(14).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Car configuration (Car #"+c.car+")"),
						c.teamName.View(),
					),
				),
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Max speed"),
						c.onTrackSpeed.View(),
						c.inPitSpeed.View(),
						c.minSpeed.View(),
					),
				),

				lipgloss.PlaceHorizontal(78, lipgloss.Center,
					lipgloss.JoinHorizontal(lipgloss.Center, ok, cancel),
				),
			),
		),
	)

	// return b.String()
}
