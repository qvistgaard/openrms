package elements

import "github.com/charmbracelet/lipgloss"

func Shortcut(shortcut string, description string) string {
	return lipgloss.JoinHorizontal(lipgloss.Right,
		lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("["),
		lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).Render(shortcut),
		lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("]"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Width(16).Render(" "+description),
	)

}
