package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/plugins/leaderboard"
	"github.com/qvistgaard/openrms/internal/plugins/rules/fuel"
	"github.com/qvistgaard/openrms/internal/plugins/rules/limbmode"
	"github.com/qvistgaard/openrms/internal/plugins/rules/pit"
	"github.com/qvistgaard/openrms/internal/state/rules"
	log "github.com/sirupsen/logrus"
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

	lb := leaderboard.New()
	ctx.Rules = rules.NewRuleList()
	ctx.Rules.Append(lb)
	ctx.Leaderboard = lb

	for _, r := range c.Rules {
		if r.Enabled {
			switch r.Plugin {
			case "fuel":
				ctx.Rules.Append(fuel.CreateFromConfig(ctx.Config, ctx.Rules))
				log.Info("fuel plugin loaded")

			case "limb-mode":
				ctx.Rules.Append(limbmode.CreateFromConfig(ctx.Config))
				log.Info("limb-mode plugin loaded")

				/*					case "damage":
									ctx.Rules.Append(&damage.Rule{})*/
			case "pit":
				ctx.Rules.Append(pit.CreatePitRule(ctx.Rules))
				log.Info("pit plugin loaded")
				/*			case "tirewear":
								ctx.Rules.Append(&tirewear.Rule{})
							case "leaderboard":
								ctx.Rules.Append(&leaderboard.Rule{})
							case "race":
								ctx.Rules.Append(race.CreateFromConfig(rm.Rules[k]))*/
			default:
				return errors.New("Unknown rule: " + r.Plugin)
			}
		}
	}
	return nil
}
