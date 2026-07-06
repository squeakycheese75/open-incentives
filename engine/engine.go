package engine

import (
	"context"
	"fmt"
	"reflect"
)

type engineImpl struct {
}

func New() Engine {
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

func evaluateCondition(c Condition, facts map[string]any) (any, bool) {
	actual, ok := facts[c.Fact]
	if !ok {
		return nil, false
	}

	switch c.Operator {
	case "eq":
		return actual, reflect.DeepEqual(actual, c.Value)

	case "neq":
		return actual, !reflect.DeepEqual(actual, c.Value)

	case "gt":
		left, right, ok := numbers(actual, c.Value)
		return actual, ok && left > right

	case "gte":
		left, right, ok := numbers(actual, c.Value)
		return actual, ok && left >= right

	case "lt":
		left, right, ok := numbers(actual, c.Value)
		return actual, ok && left < right

	case "lte":
		left, right, ok := numbers(actual, c.Value)
		return actual, ok && left <= right

	default:
		return actual, false
	}
}
