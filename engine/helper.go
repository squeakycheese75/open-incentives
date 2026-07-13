package engine

func numbers(left, right any) (float64, float64, bool) {
	l, ok := number(left)
	if !ok {
		return 0, 0, false
	}

	r, ok := number(right)
	if !ok {
		return 0, 0, false
	}

	return l, r, true
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

func isSupportedOperator(op string) bool {
	switch op {
	case "eq", "neq", "gt", "gte", "lt", "lte":
		return true
	default:
		return false
	}
}

func isSupportedActionType(actionType string) bool {
	switch actionType {
	case "percentage_discount":
		return true
	default:
		return false
	}
}

func number(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
