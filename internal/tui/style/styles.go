package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	InactiveForegroundColor = lipgloss.Color("240")

	DialogBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("24")).
			Padding(1, 0).
			MarginTop(3).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	Container = lipgloss.NewStyle().
			PaddingBottom(2).PaddingLeft(2)

	Heading = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("240")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	Button = struct {
		Focused lipgloss.Style
		Blurred lipgloss.Style
		Caption struct {
			Focused lipgloss.Style
			Blurred lipgloss.Style
		}
	}{
		Focused: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(1).PaddingRight(2),
		Blurred: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(1).PaddingRight(2),
		Caption: struct {
			Focused lipgloss.Style
			Blurred lipgloss.Style
		}{
			Focused: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(1).PaddingRight(2),
			Blurred: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(1).PaddingRight(2),
		},
	}
)