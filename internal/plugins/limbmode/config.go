package limbmode

import (
	"github.com/qvistgaard/openrms/internal/types"
)

type LimbModeConfig struct {
	MaxSpeed *uint8 `mapstructure:"max-speed"`
}

type CarSettings struct {
	Id       *types.CarId    `mapstructure:"id"`
	LimbMode *LimbModeConfig `mapstructure:"limb-mode"`
}

type Config struct {
	Car struct {
		Defaults *CarSettings   `mapstructure:"defaults"`
		Cars     []*CarSettings `mapstructure:"cars"`
	}
	Plugin struct {
		LimbMode struct {
			Enabled    bool `mapstructure:"enabled"`
			Commentary bool `mapstructure:"announcements" default:"true"`
		} `mapstructure:"limb-mode"`
	} `mapstructure:"plugins"`
}
