package engine

import (
	"reflect"
	"testing"
)

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		name       string
		condition  Condition
		facts      map[string]any
		wantActual any
		wantMatch  bool
	}{
		{
			name: "eq string matches",
			condition: Condition{
				Fact:     "customer.tier",
				Operator: "eq",
				Value:    "gold",
			},
			facts: map[string]any{
				"customer.tier": "gold",
			},
			wantActual: "gold",
			wantMatch:  true,
		},
		{
			name: "eq string does not match",
			condition: Condition{
				Fact:     "customer.tier",
				Operator: "eq",
				Value:    "gold",
			},
			facts: map[string]any{
				"customer.tier": "silver",
			},
			wantActual: "silver",
			wantMatch:  false,
		},
		{
			name: "missing fact returns false",
			condition: Condition{
				Fact:     "customer.tier",
				Operator: "eq",
				Value:    "gold",
			},
			facts:      map[string]any{},
			wantActual: nil,
			wantMatch:  false,
		},
		{
			name: "gte number matches",
			condition: Condition{
				Fact:     "cart.total",
				Operator: "gte",
				Value:    50,
			},
			facts: map[string]any{
				"cart.total": 72,
			},
			wantActual: 72,
			wantMatch:  true,
		},
		{
			name: "gte number does not match",
			condition: Condition{
				Fact:     "cart.total",
				Operator: "gte",
				Value:    50,
			},
			facts: map[string]any{
				"cart.total": 49,
			},
			wantActual: 49,
			wantMatch:  false,
		},
		{
			name: "unknown operator returns false",
			condition: Condition{
				Fact:     "cart.total",
				Operator: "wat",
				Value:    50,
			},
			facts: map[string]any{
				"cart.total": 72,
			},
			wantActual: 72,
			wantMatch:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, matched := evaluateLeafCondition(tt.condition, tt.facts)

			if !reflect.DeepEqual(actual, tt.wantActual) {
				t.Fatalf("actual = %v, want %v", actual, tt.wantActual)
			}

			if matched != tt.wantMatch {
				t.Fatalf("matched = %v, want %v", matched, tt.wantMatch)
			}
		})
	}
}

func TestEvaluateCondition_AllAny(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		facts     map[string]any
		wantMatch bool
	}{
		{
			name: "all matches when all children match",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
				},
			},
			facts: map[string]any{
				"cart.total":    72,
				"customer.tier": "gold",
			},
			wantMatch: true,
		},
		{
			name: "all fails when one child fails",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
				},
			},
			facts: map[string]any{
				"cart.total":    72,
				"customer.tier": "silver",
			},
			wantMatch: false,
		},
		{
			name: "any matches when one child matches",
			condition: Condition{
				Any: []Condition{
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
					{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
				},
			},
			facts: map[string]any{
				"customer.tier": "platinum",
			},
			wantMatch: true,
		},
		{
			name: "any fails when no children match",
			condition: Condition{
				Any: []Condition{
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
					{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
				},
			},
			facts: map[string]any{
				"customer.tier": "silver",
			},
			wantMatch: false,
		},
		{
			name: "nested any inside all matches",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{
						Any: []Condition{
							{Fact: "customer.tier", Operator: "eq", Value: "gold"},
							{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
						},
					},
				},
			},
			facts: map[string]any{
				"cart.total":    72,
				"customer.tier": "gold",
			},
			wantMatch: true,
		},
		{
			name: "nested any inside all fails",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{
						Any: []Condition{
							{Fact: "customer.tier", Operator: "eq", Value: "gold"},
							{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
						},
					},
				},
			},
			facts: map[string]any{
				"cart.total":    72,
				"customer.tier": "silver",
			},
			wantMatch: false,
		},
		{
			name: "all fails when child fact is missing",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
				},
			},
			facts: map[string]any{
				"cart.total": 72,
			},
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, matched := evaluateCondition(tt.condition, tt.facts)

			if matched != tt.wantMatch {
				t.Fatalf("matched = %v, want %v", matched, tt.wantMatch)
			}
		})
	}
}
