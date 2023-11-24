package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/tui/style"
	"strconv"
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
	driverName    textinput.Model
	minSpeed      textinput.Model
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(1).PaddingRight(2)
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(1).PaddingRight(2)
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func InitialCarConfigurationModel() CarConfiguration {
	m := CarConfiguration{}

	m.driverName = textinput.New()
	m.driverName.Focus()
	m.driverName.PromptStyle = focusedStyle.Copy().Width(18)
	m.driverName.Prompt = "Driver name:"
	m.driverName.TextStyle = focusedStyle
	m.driverName.CharLimit = 64
	m.driverName.Cursor.SetMode(cursor.CursorBlink)

	m.onTrackSpeed = textinput.New()
	m.onTrackSpeed.PromptStyle = blurredStyle.Copy().Width(18)
	m.onTrackSpeed.Placeholder = "100"
	m.onTrackSpeed.Prompt = "On track:"
	m.onTrackSpeed.TextStyle = focusedStyle.Copy().Underline(true)
	m.onTrackSpeed.CharLimit = 3

	m.inPitSpeed = textinput.New()
	m.inPitSpeed.Placeholder = "75"
	m.inPitSpeed.PromptStyle = blurredStyle.Copy().Width(18)
	m.inPitSpeed.Prompt = "In pit lane:"
	m.inPitSpeed.TextStyle = focusedStyle
	m.inPitSpeed.CharLimit = 3

	m.minSpeed = textinput.New()
	m.minSpeed.Placeholder = "75"
	m.minSpeed.PromptStyle = blurredStyle.Copy().Width(18)
	m.minSpeed.Prompt = "Min speed:"
	m.minSpeed.TextStyle = focusedStyle
	m.minSpeed.CharLimit = 3

	m.inputs = []*textinput.Model{&m.driverName, &m.onTrackSpeed, &m.inPitSpeed, &m.minSpeed}
	m.focusIndexMax = len(m.inputs) + 1
	return m
}

func (c CarConfiguration) Init() tea.Cmd {
	return textinput.Blink
}

func (c CarConfiguration) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	c.inputs = []*textinput.Model{&c.driverName, &c.onTrackSpeed, &c.inPitSpeed, &c.minSpeed}

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
						DriverName:  c.driverName.Value(),
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
				(*input).PromptStyle = blurredStyle.Copy().Width(18)
			}

			if c.focusIndex >= 0 && c.focusIndex < len(c.inputs) {
				(c.inputs[c.focusIndex]).PromptStyle = focusedStyle.Copy().Width(18)
				cmds[c.focusIndex] = (c.inputs[c.focusIndex]).Focus()
			}
			return c, tea.Batch(cmds...)
		}
	case commands.OpenCarConfiguration:
		c.focusIndex = 0
		c.car = msg.CarId
		c.driverName.SetValue(msg.DriverName)
		c.inPitSpeed.SetValue(msg.MaxPitSpeed)
		c.onTrackSpeed.SetValue(msg.MaxSpeed)
		c.minSpeed.SetValue(msg.MinSpeed)
		c.driverName.Focus()
		return c, nil
	}

	var cmd tea.Cmd
	if c.focusIndex >= 0 && c.focusIndex < len(c.inputs) {
		*c.inputs[c.focusIndex], cmd = (*c.inputs[c.focusIndex]).Update(msg)
	}
	return c, cmd
}

func (c CarConfiguration) View() string {
	var ok string
	if c.focusIndex == len(c.inputs) {
		ok = focusedStyle.Copy().Render("[ Save ]")
	} else {
		ok = blurredStyle.Copy().Render("[ Save ]")
	}

	var cancel string

	if c.focusIndex == len(c.inputs)+1 {
		cancel = focusedStyle.Copy().Render("[ Cancel ]")
	} else {
		cancel = blurredStyle.Copy().Render("[ Cancel ]")
	}

	return lipgloss.PlaceHorizontal(200,
		lipgloss.Center,

		style.DialogBox.Width(77).Height(14).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				style.Container.Render(
					lipgloss.JoinVertical(lipgloss.Top,
						style.Heading.Width(72).Render("Car configuration (Car #"+c.car+")"),
						c.driverName.View(),
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
					lipgloss.JoinHorizontal(lipgloss.Center, ok, cancel, strconv.Itoa(c.focusIndex)),
				),
			),
		),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	// return b.String()
}