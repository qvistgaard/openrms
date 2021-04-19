package state

import "sort"

type Rule interface {
	InitializeCarState(car *Car)
	InitializeRaceState(race *Course)
}

type PitRule interface {
	HandlePitStop(car *Car)
	Priority() uint8
}

type Rules interface {
	All() []Rule
	PitRules() []PitRule
	Append(rule Rule)
}

type RuleList struct {
	rules    []Rule
	pitRules []PitRule
}

func CreateRuleList() Rules {
	rl := new(RuleList)
	rl.rules = []Rule{}
	rl.pitRules = []PitRule{}
	return rl
}

func (r *RuleList) All() []Rule {
	return r.rules
}

func (r *RuleList) PitRules() []PitRule {
	return r.pitRules
}

func (r *RuleList) Append(rule Rule) {
	r.rules = append(r.rules, rule)
	if pr, ok := rule.(PitRule); ok {
		r.pitRules = append(r.pitRules, pr)
	}
	sort.Slice(r.pitRules, func(i, j int) bool {
		return r.pitRules[i].Priority() > r.pitRules[j].Priority()
	})
}
