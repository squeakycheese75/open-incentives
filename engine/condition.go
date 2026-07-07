package engine

import "reflect"

func evaluateCondition(c Condition, facts map[string]any) (any, bool) {
	if len(c.All) > 0 {
		for _, child := range c.All {
			_, matched := evaluateCondition(child, facts)
			if !matched {
				return nil, false
			}
		}
		return nil, true
	}

	if len(c.Any) > 0 {
		for _, child := range c.Any {
			_, matched := evaluateCondition(child, facts)
			if matched {
				return nil, true
			}
		}
		return nil, false
	}

	return evaluateLeafCondition(c, facts)
}

func evaluateLeafCondition(c Condition, facts map[string]any) (any, bool) {
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
