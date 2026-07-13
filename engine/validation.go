package engine

import (
	"errors"
	"fmt"
)

func (r Rule) Validate() error {
	if r.ID == "" {
		return errors.New("rule id is required")
	}

	if err := r.Conditions.Validate(); err != nil {
		return fmt.Errorf("conditions: %w", err)
	}

	if len(r.Actions) == 0 {
		return errors.New("at least one action is required")
	}

	for i, action := range r.Actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("actions[%d]: %w", i, err)
		}
	}

	return nil
}

func (c Condition) Validate() error {
	hasAll := len(c.All) > 0
	hasAny := len(c.Any) > 0
	hasLeaf := c.Fact != "" || c.Operator != "" || c.Value != nil

	count := 0
	if hasAll {
		count++
	}
	if hasAny {
		count++
	}
	if hasLeaf {
		count++
	}

	if count != 1 {
		return errors.New("condition must be exactly one of all, any, or leaf")
	}

	if hasAll {
		for i, child := range c.All {
			if err := child.Validate(); err != nil {
				return fmt.Errorf("all[%d]: %w", i, err)
			}
		}
		return nil
	}

	if hasAny {
		for i, child := range c.Any {
			if err := child.Validate(); err != nil {
				return fmt.Errorf("any[%d]: %w", i, err)
			}
		}
		return nil
	}

	if c.Fact == "" {
		return errors.New("fact is required")
	}

	if c.Operator == "" {
		return errors.New("operator is required")
	}

	if !isSupportedOperator(c.Operator) {
		return fmt.Errorf("unsupported operator %q", c.Operator)
	}

	return nil
}

func (a Action) Validate() error {
	if a.Type == "" {
		return errors.New("action type is required")
	}

	switch a.Type {
	case "percentage_discount":
		value, ok := a.Params["value"]
		if !ok {
			return errors.New("percentage_discount requires param value")
		}

		n, ok := number(value)
		if !ok {
			return errors.New("percentage_discount value must be a number")
		}

		if n <= 0 || n > 100 {
			return errors.New("percentage_discount value must be greater than 0 and less than or equal to 100")
		}

		return nil

	default:
		return fmt.Errorf("unsupported action type %q", a.Type)
	}
}
