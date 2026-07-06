package engine

func numbers(left any, right any) (float64, float64, bool) {
	l, ok := toFloat64(left)
	if !ok {
		return 0, 0, false
	}

	r, ok := toFloat64(right)
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
