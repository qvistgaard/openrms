package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/plugins/rules/damage"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/rules/leaderboard"
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/plugins/rules/tirewear"
	"github.com/qvistgaard/openrms/internal/state"
)

type RuleConfig struct {
	Rules []struct {
		Plugin string
	}
}

func CreateRules(ctx *context.Context) error {
	c := &RuleConfig{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return err
	}
	ctx.Rules = state.CreateRuleList()
	for _, r := range c.Rules {
		switch r.Plugin {
		case "fuel":
			ctx.Rules.Append(&fuel.Consumption{})
		case "limb-mode":
			ctx.Rules.Append(&limbmode.LimbMode{})
		case "damage":
			ctx.Rules.Append(&damage.Damage{})
		case "pit":
			ctx.Rules.Append(pit.CreatePitRule(ctx))
		case "tirewear":
			ctx.Rules.Append(&tirewear.TireWear{})
		case "leaderboard":
			ctx.Rules.Append(&leaderboard.Rule{})
		default:
			return errors.New("Unknown rule: " + r.Plugin)
		}
	}
	return nil
}
