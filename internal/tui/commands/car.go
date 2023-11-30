package commands

type SetCarDriverName string

type OpenCarConfiguration struct {
	CarId       string
	MaxSpeed    string
	MaxPitSpeed string
	DriverName  string
	MinSpeed    string
}

type SaveCarConfiguration OpenCarConfiguration

/*
func OpenCarConfiguration(id string, maxSpeed string, maxPitSpeed string, driverName string) tea.Msg {
	return CarConfigurationCommand{
		CarId:       id,
		SetMaxSpeed:    maxSpeed,
		MaxPitSpeed: maxPitSpeed,
		DriverName:  driverName,
	}

}
*/
