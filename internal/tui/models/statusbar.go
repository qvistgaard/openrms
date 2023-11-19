package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

type StatusBar struct {
	width int
}

func (s StatusBar) Init() tea.Cmd {
	return nil
}

func (s StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.(tea.WindowSizeMsg).Width
	}
	return s, nil
}

func (s StatusBar) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("24")).
		PaddingLeft(1).
		Width(s.width).
		Render("Status: XXX", strconv.Itoa(s.width))
}
