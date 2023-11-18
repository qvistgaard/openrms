package style

import "github.com/charmbracelet/lipgloss"

var (
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
)
