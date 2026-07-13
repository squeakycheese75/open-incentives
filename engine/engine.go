package engine

import (
	"context"
	"encoding/json"
	"fmt"
)

type Runtime interface {
	ValidateRules(raw json.RawMessage) error
	Evaluate(ctx context.Context, req EvaluationRequest) (EvaluationResult, error)
}

type engineImpl struct {
}

func New() Runtime {
	return &engineImpl{}
}

func (e *engineImpl) Evaluate(ctx context.Context, req EvaluationRequest) (EvaluationResult, error) {
	var result EvaluationResult

	for _, rule := range req.Rules {
		matched := true

		for _, condition := range rule.Conditions.All {
			actual, passed := evaluateCondition(condition, req.Facts)

			result.Trace = append(result.Trace, TraceEntry{
				RuleID:   rule.ID,
				Fact:     condition.Fact,
				Operator: condition.Operator,
				Expected: condition.Value,
				Actual:   actual,
				Passed:   passed,
				Message:  fmt.Sprintf("%s %s %v", condition.Fact, condition.Operator, condition.Value),
			})

			if !passed {
				matched = false
			}
		}

		if matched {
			result.MatchedRules = append(result.MatchedRules, rule.ID)
			result.Actions = append(result.Actions, rule.Actions...)
		} else {
			result.RejectedRules = append(result.RejectedRules, rule.ID)
		}
	}

	return result, nil
}

func (e *engineImpl) ValidateRules(raw json.RawMessage) error {
	var rules []Rule

	if err := json.Unmarshal(raw, &rules); err != nil {
		return fmt.Errorf("decode rules: %w", err)
	}

	if len(rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}

	for i, rule := range rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("rule %d: %w", i, err)
		}
	}

	return nil
}
