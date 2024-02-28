package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/plugins/telemetry"
	"github.com/qvistgaard/openrms/internal/tui/commands"
	table "github.com/qvistgaard/openrms/internal/tui/elements"
	"github.com/qvistgaard/openrms/internal/tui/messages"
	"github.com/qvistgaard/openrms/internal/types"
	"math"
	"strconv"
	"strings"
	"time"
)

type Leaderboard struct {
	table         table.Model
	width         int
	rows          []table.Row
	rowsId        []types.CarId
	raceTelemetry telemetry.Race
	colors        map[uint]string
}

const (
	DESLOTTED = "⛐"
	LIMBMODE  = "⛍"
	PIT       = "P" // "🛠"
)

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
		{Title: alignRight.Width(3).Align(lipgloss.Center).Render("P"), Width: 3},
		{Title: "Name", Width: 120},
		{Title: alignRight.Width(3).Render("#"), Width: 3},
		{Title: alignRight.Width(4).Render("Fuel"), Width: 4},
		{Title: lipgloss.NewStyle().Width(7).Align(lipgloss.Right).Render("Lap"), Width: 7},
		{Title: alignRight.Width(7).Align(lipgloss.Right).Render("Delta"), Width: 7},
		{Title: alignRight.Width(7).Align(lipgloss.Right).Render("Best"), Width: 7},
		{Title: alignRight.Width(4).Render("Laps"), Width: 4},
		{Title: alignRight.Width(3).Align(lipgloss.Center).Render("Pit"), Width: 3},
		{Title: alignRight.Width(3).Align(lipgloss.Center).Render(LIMBMODE), Width: 3},
		{Title: " " + DESLOTTED, Width: 3},
		{Title: " MS", Width: 3},
	}
)

func InitialLeaderboardModel() Leaderboard {
	l := Leaderboard{}
	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle

	s.RenderCell = func(model table.Model, value string, position table.CellPosition) string {
		if position.RowID == 0 && position.Column == 0 {
			return s.Cell.Copy().Background(lipgloss.Color("250")).Foreground(lipgloss.Color("16")).Bold(true).Render(value)
		}
		if strings.Contains(value, DESLOTTED) {
			return s.Cell.Copy().Background(lipgloss.Color("124")).Foreground(lipgloss.Color("15")).Bold(true).Render(value)
		}
		if strings.Contains(value, LIMBMODE) {
			return s.Cell.Copy().Background(lipgloss.Color("220")).Foreground(lipgloss.Color("16")).Bold(true).Render(value)
		}
		if position.Column == 8 && len(strings.Trim(value, " ")) > 0 {
			return s.Cell.Copy().Background(lipgloss.Color("18")).Foreground(lipgloss.Color("255")).Bold(true).Render(value)
		}

		if position.Column == 3 {
			parseInt, err := strconv.ParseInt(value, 10, 32)
			if err == nil {
				if parseInt < 25 && parseInt > 15 {
					return s.Cell.Copy().Background(lipgloss.Color("220")).Foreground(lipgloss.Color("16")).Bold(true).Render(value)
				}
				if parseInt > 15 {
					return s.Cell.Copy().Background(lipgloss.Color("124")).Foreground(lipgloss.Color("15")).Bold(true).Render(value)
				}
			}
		}

		if position.Column == 5 && value[0] == ("-")[0] {
			return s.Cell.Copy().Foreground(lipgloss.Color("40")).Render(value)
		}

		if position.Column == 2 {
			v, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 32)
			c, ok := l.colors[uint(v)]
			if ok {
				var bg, fg string
				switch c {
				case "red":
					bg = "124"
					fg = "255"
				case "orange":
					bg = "208"
					fg = "232"
				case "blue":
					bg = "19"
					fg = "255"
				case "black":
					bg = "16"
					fg = "254"
				case "green":
					bg = "34"
					fg = "16"
				case "white":
					bg = "253"
					fg = "16"
				case "yellow":
					bg = "214"
					fg = "16"
				case "purple":
					bg = "91"
					fg = "255"
				}

				if bg != "" && fg != "" {
					return s.Cell.Copy().Foreground(lipgloss.Color(fg)).Background(lipgloss.Color(bg)).Render(value)
				}
			}
		}

		if position.IsRowSelected {
			return s.Cell.Copy().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("24")).
				Render(value)
		}

		return s.Cell.Copy().Render(value)
	}
	l.colors = make(map[uint]string)
	l.table = table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithStyles(s))
	return l
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
				carIdString := strconv.Itoa(int(l.rowsId[l.table.Cursor()]))
				carId, _ := types.IdFromString(carIdString)
				return commands.OpenCarConfiguration{
					CarId:       carIdString,
					MaxSpeed:    fmt.Sprintf("%d", l.raceTelemetry[carId].MaxSpeed),
					MaxPitSpeed: fmt.Sprintf("%d", l.raceTelemetry[carId].MaxPitSpeed),
					MinSpeed:    fmt.Sprintf("%d", l.raceTelemetry[carId].MinSpeed),
					DriverName:  l.rows[l.table.Cursor()][1],
				}
			}
		case "e":
			return l, func() tea.Msg {
				carIdString := strconv.Itoa(int(l.rowsId[l.table.Cursor()]))
				return commands.ToggleEnableDisableCar{
					CarId: carIdString,
				}
			}
		}
		l.table, cmd = l.table.Update(msg)
		break

	case tea.WindowSizeMsg:
		l.width = msg.(tea.WindowSizeMsg).Width
		columns[1] = table.Column{Title: "Name", Width: l.width - 71}

		l.table.SetWidth(l.width)
		l.table.SetHeight(msg.(tea.WindowSizeMsg).Height - 9)
		l.table.SetColumns(columns)

	case messages.Update:
		l.rows = make([]table.Row, 0)
		l.rowsId = make([]types.CarId, 0)
		raceTelemetry := msg.(messages.Update).RaceTelemetry
		for k, v := range raceTelemetry.Sort() {
			var inPitString string
			var lmMode string
			var deslotted string
			if v.InPit {
				if v.PitStopActive {
					inPitString = strconv.Itoa(int(v.PitStopSequence))
				} else {
					inPitString = PIT
				}
			} else {
				inPitString = ""
			}
			if v.LimbMode {
				lmMode = LIMBMODE
			} else {
				lmMode = ""
			}
			if v.Deslotted {
				deslotted = DESLOTTED
			} else {
				deslotted = ""
			}

			var team string
			if v.Enabled {
				team = v.Team
			} else {
				team = v.Team + " (Disabled)"
			}

			l.colors[v.Number] = v.Color

			l.rowsId = append(l.rowsId, v.Id)
			l.rows = append(l.rows, table.Row{

				alignRight.AlignHorizontal(lipgloss.Center).Render(strconv.Itoa(k + 1)),
				team,
				alignRight.Width(3).AlignHorizontal(lipgloss.Right).Render(strconv.Itoa(int(v.Number))),
				alignRight.Width(4).AlignHorizontal(lipgloss.Right).Render(fmt.Sprintf("%.f", v.Fuel)),
				alignRight.Width(7).AlignHorizontal(lipgloss.Right).Render(formatDurationSecondsMilliseconds(v.Last.Time)),
				alignRight.Width(7).AlignHorizontal(lipgloss.Right).Render(formatDurationSecondsMilliseconds(v.Delta)),
				alignRight.Width(7).AlignHorizontal(lipgloss.Right).Render(formatDurationSecondsMilliseconds(v.Best)),
				alignRight.Width(4).Render(strconv.Itoa(int(v.Last.Number))),
				" " + inPitString,
				" " + lmMode,
				" " + deslotted,
				alignRight.Width(3).AlignHorizontal(lipgloss.Right).Render(strconv.Itoa(int(v.MaxSpeed))),
			})
		}
		l.table.SetRows(l.rows)
		l.raceTelemetry = raceTelemetry
	}

	return l, cmd

}

func (l Leaderboard) View() string {
	return l.table.View()

}

func formatDurationSecondsMilliseconds(d time.Duration) string {
	// Extract seconds and milliseconds
	seconds := math.Floor(d.Seconds())

	milliseconds := float64(d.Milliseconds()) - (seconds * 1000)

	// Format as "ss.ms"
	return fmt.Sprintf("%.0f.%03.0fs", seconds, milliseconds)
}
