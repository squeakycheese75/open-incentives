package adapters

import (
	engine "github.com/squeakycheese75/open-incentives-engine"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func mapResult(
	result engine.EvaluationResult,
	rules []engine.Rule,
) domain.RuleEvaluationResult {
	rulesByID := make(map[string]engine.Rule, len(rules))

	for _, rule := range rules {
		rulesByID[rule.ID] = rule
	}

	matches := make([]domain.MatchedRule, 0, len(result.MatchedRules))

	for _, ruleID := range result.MatchedRules {
		rule := rulesByID[ruleID]

		actions := make([]domain.Action, 0, len(rule.Actions))
		for _, action := range rule.Actions {
			actions = append(actions, domain.Action{
				Type:   action.Type,
				Params: action.Params,
			})
		}

		matches = append(matches, domain.MatchedRule{
			RuleID:  ruleID,
			Actions: actions,
		})
	}

	return domain.RuleEvaluationResult{
		MatchedRules: matches,
	}
}
