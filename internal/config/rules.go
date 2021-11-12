package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/state/rx/rules"
)

type RuleConfig struct {
	Rules []struct {
		Plugin  string
		Enabled bool
	}
}

type RuleConfigMaps struct {
	Rules []map[string]interface{}
}

func CreateRules(ctx *application.Context) error {
	c := &RuleConfig{}
	rm := &RuleConfigMaps{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return err
	}
	err = mapstructure.Decode(ctx.Config, rm)
	if err != nil {
		return err
	}

	ctx.Rules = rules.NewRuleList()
	for k, r := range c.Rules {
		if r.Enabled {
			switch r.Plugin {
			case "fuel":
				ctx.Rules.Append(fuel.Create(rm.Rules[k], ctx.Postprocessors))
				/*
					case "limb-mode":
						ctx.Rules.Append(limbmode.CreateFromConfig(rm.Rules[k]))
					case "damage":
						ctx.Rules.Append(&damage.Rule{})*/
			case "pit":
				ctx.Rules.Append(pit.CreatePitRule(ctx.Rules))
				/*			case "tirewear":
								ctx.Rules.Append(&tirewear.Rule{})
							case "leaderboard":
								ctx.Rules.Append(&leaderboard.Rule{})
							case "race":
								ctx.Rules.Append(race.Create(rm.Rules[k]))*/
			default:
				return errors.New("Unknown rule: " + r.Plugin)
			}
		}
	}
	return nil
}
