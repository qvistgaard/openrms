package car

import "github.com/qvistgaard/openrms/internal/types"

type PitLaneConfig struct {
	MaxSpeed *uint8 `mapstructure:"max-speed"`
}

type Settings struct {
	Id           *types.CarId   `mapstructure:"id"`
	MaxSpeed     *uint8         `mapstructure:"max-speed"`
	PitLane      *PitLaneConfig `mapstructure:"pit-lane"`
	MaxBreaking  *uint8         `mapstructure:"max-breaking"`
	MinSpeed     *uint8         `mapstructure:"min-speed"`
	Drivers      *types.Drivers `mapstructure:"drivers"`
	Team         *string        `mapstructure:"team"`
	Number       *uint          `mapstructure:"number"`
	Manufacturer *string        `mapstructure:"manufacturer"`
	Color        *string        `mapstructure:"color"`
}

type Config struct {
	Car struct {
		Defaults *Settings   `mapstructure:"defaults"`
		Cars     []*Settings `mapstructure:"cars"`
	}
}
