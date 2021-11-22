package rules

import (
	"github.com/qvistgaard/openrms/internal/state/rules/car"
	"github.com/qvistgaard/openrms/internal/state/rules/pit"
	"github.com/qvistgaard/openrms/internal/state/rules/race"
	"sort"
)

type Rule interface {
	Priority() int
	Name() string
}

type Rules interface {
	CarRules() []car.Rule
	CarRule(name string) car.Rule
	RaceRules() []race.Rule
	PitRules() []pit.Rule
	Append(Rule)
}

type CarRule struct {
	name string
	rule car.Rule
}

type RuleList struct {
	carRules  []car.Rule
	pitRules  []pit.Rule
	raceRules []race.Rule
}

func NewRuleList() *RuleList {
	return &RuleList{
		carRules:  []car.Rule{},
		pitRules:  []pit.Rule{},
		raceRules: []race.Rule{},
	}
}

func (r *RuleList) CarRule(name string) car.Rule {
	for _, v := range r.carRules {
		if v.(Rule).Name() == name {
			return v
		}
	}
	return nil
}

func (r *RuleList) CarRules() []car.Rule {
	return r.carRules
}

func (r *RuleList) PitRules() []pit.Rule {
	return r.pitRules
}

func (r *RuleList) RaceRules() []race.Rule {
	return r.raceRules
}

func (r *RuleList) Append(rule Rule) {
	if rule, ok := rule.(car.Rule); ok {
		r.carRules = append(r.carRules, rule)
	}
	if rule, ok := rule.(pit.Rule); ok {
		r.pitRules = append(r.pitRules, rule)
	}
	if rule, ok := rule.(race.Rule); ok {
		r.raceRules = append(r.raceRules, rule)
	}

	sort.Slice(r.carRules, func(i, j int) bool {
		return r.carRules[i].(Rule).Priority() < r.carRules[j].(Rule).Priority()
	})
	sort.Slice(r.pitRules, func(i, j int) bool {
		return r.pitRules[i].(Rule).Priority() < r.pitRules[j].(Rule).Priority()
	})
	sort.Slice(r.raceRules, func(i, j int) bool {
		return r.raceRules[i].(Rule).Priority() < r.raceRules[j].(Rule).Priority()
	})
}
