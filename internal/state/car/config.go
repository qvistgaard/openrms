package car

import "github.com/qvistgaard/openrms/internal/types"

type PitLaneConfig struct {
	MaxSpeed *uint8 `mapstructure:"max-speed"`
}

type CarSettings struct {
	Id          *types.Id      `mapstructure:"id"`
	MaxSpeed    *uint8         `mapstructure:"max-speed"`
	PitLane     *PitLaneConfig `mapstructure:"pit-lane"`
	MaxBreaking *uint8         `mapstructure:"max-breaking"`
	MinSpeed    *uint8         `mapstructure:"min-speed"`
	Drivers     *types.Drivers `mapstructure:"drivers"`
	Team        *string        `mapstructure:"team"`
}

type Config struct {
	Car struct {
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
}
