package car

import "github.com/qvistgaard/openrms/internal/types"

type PitLaneConfig struct {
	MaxSpeed *types.Percent `mapstructure:"max-speed"`
}

type CarSettings struct {
	Id          *types.Id      `mapstructure:"id"`
	MaxSpeed    *types.Percent `mapstructure:"max-speed"`
	PitLane     *PitLaneConfig `mapstructure:"pit-lane"`
	MaxBreaking *types.Percent `mapstructure:"max-breaking"`
	MinSpeed    *types.Percent `mapstructure:"min-speed"`
}

type Config struct {
	Car struct {
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
}
