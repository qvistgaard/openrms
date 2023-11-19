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
	alignRight  = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right)
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
		{Title: alignRight.Width(3).Render("P"), Width: 3},
		{Title: "Name", Width: 120 - 60},
		{Title: alignRight.Width(3).Render("#"), Width: 3},
		{Title: alignRight.Width(4).Render("Fuel"), Width: 4},
		{Title: alignRight.Width(7).Render("Lap"), Width: 7},
		{Title: alignRight.Width(7).Render("Delta"), Width: 7},
		{Title: alignRight.Width(7).Render("Best"), Width: 7},
		{Title: alignRight.Width(4).Render("Laps"), Width: 4},
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
				carIdString := l.rows[l.table.Cursor()][2]
				carId, _ := types.IdFromString(carIdString)
				return commands.OpenCarConfiguration{
					CarId:       carIdString,
					MaxSpeed:    fmt.Sprintf("%.f", l.raceTelemetry[carId].MaxSpeed),
					MaxPitSpeed: fmt.Sprintf("%.f", l.raceTelemetry[carId].MaxPitSpeed),
					MinSpeed:    fmt.Sprintf("%.f", l.raceTelemetry[carId].MinSpeed),
					DriverName:  l.rows[l.table.Cursor()][1],
				}
			}

		}
		l.table, cmd = l.table.Update(msg)
		break

	case tea.WindowSizeMsg:
		l.width = msg.(tea.WindowSizeMsg).Width
		columns[1] = table.Column{Title: "Name", Width: l.width - 66}

		l.table.SetWidth(l.width)
		l.table.SetColumns(columns)
	case types.RaceTelemetry:
		l.rows = make([]table.Row, 0)
		for k, v := range msg.(types.RaceTelemetry).Sort() {
			l.rows = append(l.rows, table.Row{
				alignRight.Width(3).Render(strconv.Itoa(k + 1)),
				v.Name,
				alignRight.Width(3).Render(strconv.Itoa(int(v.Car))),
				alignRight.Width(4).Render(fmt.Sprintf("%.f", v.Fuel)),
				alignRight.Width(7).Render(v.Last.String()),
				alignRight.Width(7).Render(v.Delta.String()),
				alignRight.Width(7).Render(v.Best.String()),
				alignRight.Width(4).Render(strconv.Itoa(int(v.Laps.LapNumber))),
			})
		}
		l.table.SetRows(l.rows)
		l.raceTelemetry = msg.(types.RaceTelemetry)
	}

	return l, cmd

}

func (l Leaderboard) View() string {
	return l.table.View()
}
