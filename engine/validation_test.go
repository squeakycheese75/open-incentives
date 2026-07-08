package engine

import "testing"

func TestRuleValidate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: Rule{
				ID: "rule_1",
				Conditions: Condition{
					Fact:     "cart.total",
					Operator: "gte",
					Value:    50,
				},
				Actions: []Action{
					{
						Type: "percentage_discount",
						Params: map[string]any{
							"value": 10,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing id",
			rule: Rule{
				Conditions: Condition{
					Fact:     "cart.total",
					Operator: "gte",
					Value:    50,
				},
				Actions: []Action{
					{
						Type: "percentage_discount",
						Params: map[string]any{
							"value": 10,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing actions",
			rule: Rule{
				ID: "rule_1",
				Conditions: Condition{
					Fact:     "cart.total",
					Operator: "gte",
					Value:    50,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid condition",
			rule: Rule{
				ID:         "rule_1",
				Conditions: Condition{},
				Actions: []Action{
					{
						Type: "percentage_discount",
						Params: map[string]any{
							"value": 10,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid action",
			rule: Rule{
				ID: "rule_1",
				Conditions: Condition{
					Fact:     "cart.total",
					Operator: "gte",
					Value:    50,
				},
				Actions: []Action{
					{
						Type: "unknown_action",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestConditionValidate(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		wantErr   bool
	}{
		{
			name: "valid leaf condition",
			condition: Condition{
				Fact:     "cart.total",
				Operator: "gte",
				Value:    50,
			},
			wantErr: false,
		},
		{
			name: "valid all condition",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid any condition",
			condition: Condition{
				Any: []Condition{
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
					{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
				},
			},
			wantErr: false,
		},
		{
			name:      "empty condition is invalid",
			condition: Condition{},
			wantErr:   true,
		},
		{
			name: "all and any together is invalid",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
				},
				Any: []Condition{
					{Fact: "customer.tier", Operator: "eq", Value: "gold"},
				},
			},
			wantErr: true,
		},
		{
			name: "group and leaf together is invalid",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
				},
				Fact:     "customer.tier",
				Operator: "eq",
				Value:    "gold",
			},
			wantErr: true,
		},
		{
			name: "missing fact is invalid",
			condition: Condition{
				Operator: "eq",
				Value:    "gold",
			},
			wantErr: true,
		},
		{
			name: "missing operator is invalid",
			condition: Condition{
				Fact:  "customer.tier",
				Value: "gold",
			},
			wantErr: true,
		},
		{
			name: "unsupported operator is invalid",
			condition: Condition{
				Fact:     "customer.tier",
				Operator: "contains",
				Value:    "gold",
			},
			wantErr: true,
		},
		{
			name: "invalid nested child is invalid",
			condition: Condition{
				All: []Condition{
					{Fact: "cart.total", Operator: "gte", Value: 50},
					{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.condition.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestActionValidate(t *testing.T) {
	tests := []struct {
		name    string
		action  Action
		wantErr bool
	}{
		{
			name: "valid percentage discount",
			action: Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": 10,
				},
			},
			wantErr: false,
		},
		{
			name: "missing action type",
			action: Action{
				Params: map[string]any{
					"value": 10,
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported action type",
			action: Action{
				Type: "fixed_discount",
				Params: map[string]any{
					"value": 10,
				},
			},
			wantErr: true,
		},
		{
			name: "percentage discount missing value",
			action: Action{
				Type:   "percentage_discount",
				Params: map[string]any{},
			},
			wantErr: true,
		},
		{
			name: "percentage discount value is not number",
			action: Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": "10",
				},
			},
			wantErr: true,
		},
		{
			name: "percentage discount value is zero",
			action: Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": 0,
				},
			},
			wantErr: true,
		},
		{
			name: "percentage discount value above 100",
			action: Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": 101,
				},
			},
			wantErr: true,
		},
		{
			name: "percentage discount value is 100",
			action: Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": 100,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
