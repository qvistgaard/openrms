package rules

import (
	"github.com/qvistgaard/openrms/internal/state/rx/rules/car"
	"github.com/qvistgaard/openrms/internal/state/rx/rules/pit"
	"sort"
)

type Rule interface {
	Priority() int
}

type Rules interface {
	CarRules() []car.Rule
	PitRules() []pit.Rule
	Append(Rule)
}

type RuleList struct {
	carRules []car.Rule
	pitRules []pit.Rule
}

func NewRuleList() *RuleList {
	return &RuleList{
		carRules: []car.Rule{},
		pitRules: []pit.Rule{},
	}
}

/*func (r *RuleList) PitRules() []Rule {
	return r.pitRules
}
*/

func (r *RuleList) CarRules() []car.Rule {
	return r.carRules
}

func (r *RuleList) PitRules() []pit.Rule {
	return r.pitRules
}

func (r *RuleList) Append(rule Rule) {
	if rule, ok := rule.(car.Rule); ok {
		r.carRules = append(r.carRules, rule)
	}
	if rule, ok := rule.(pit.Rule); ok {
		r.pitRules = append(r.pitRules, rule)
	}

	sort.Slice(r.carRules, func(i, j int) bool {
		return r.carRules[i].(Rule).Priority() > r.carRules[j].(Rule).Priority()
	})
	sort.Slice(r.pitRules, func(i, j int) bool {
		return r.pitRules[i].(Rule).Priority() > r.pitRules[j].(Rule).Priority()
	})
}
