package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	"github.com/qvistgaard/openrms/internal/types"
	"strconv"
)

type Leaderboard struct {
	table         table.Model
	width         int
	rows          []table.Row
	raceTelemetry types.RaceTelemetry
}

var (
	borderStyle = lipgloss.NormalBorder()
	headerStyle = table.DefaultStyles().Header.
			AlignVertical(lipgloss.Right).
			BorderStyle(borderStyle).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			PaddingTop(1).
			Bold(true)

	selectedStyle = table.DefaultStyles().Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("24")).
			AlignVertical(lipgloss.Right).
			Bold(false)

	columns = []table.Column{
		{Title: "P", Width: 3},
		{Title: "Name", Width: 120 - 67},
		{Title: "#", Width: 3},
		{Title: "Fuel", Width: 10},
		{Title: "Lap", Width: 7},
		{Title: "Delta", Width: 7},
		{Title: "Best", Width: 7},
		{Title: "Laps", Width: 5},
	}
)

func NewLeaderBoard() Leaderboard {
	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle

	return Leaderboard{
		table: table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
			table.WithHeight(10),
			table.WithStyles(s),
		),
	}
}

func (l Leaderboard) Init() tea.Cmd {
	return nil
}

func (l Leaderboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "enter":
			return l, nil
		case "c":
			return l, func() tea.Msg {
				return commands.OpenCarConfiguration{
					CarId:       l.rows[l.table.Cursor()][3],
					MaxSpeed:    "100",
					MaxPitSpeed: "80",
					DriverName:  l.rows[l.table.Cursor()][1],
				}
			}

		}
		l.table, cmd = l.table.Update(msg)
		break

	case tea.WindowSizeMsg:
		l.width = msg.(tea.WindowSizeMsg).Width
		columns[1] = table.Column{Title: "Name", Width: l.width - 67}

		l.table.SetWidth(l.width)
		l.table.SetColumns(columns)
	case types.RaceTelemetry:
		l.rows = make([]table.Row, 0)
		for k, v := range msg.(types.RaceTelemetry).Sort() {
			l.rows = append(l.rows, table.Row{strconv.Itoa(k + 1), v.Name + "fdafas", strconv.Itoa(int(v.Car)), fmt.Sprintf("%f", v.Fuel), v.Last.String(), v.Delta.String(), v.Best.String(), strconv.Itoa(int(v.Laps.LapNumber))})
		}
		l.table.SetRows(l.rows)
		l.raceTelemetry = msg.(types.RaceTelemetry)
	}

	return l, cmd

}

func (l Leaderboard) View() string {
	return l.table.View()
}
