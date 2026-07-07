# Open Incentives Engine

> A lightweight, embeddable rules engine that powers Open Incentives.

The engine is intentionally generic. It does not know about:

- Promotions
- Coupons
- Customers
- Carts
- Stripe
- Discounts

It only understands:

```text
Facts + Rules -> Actions + Decisions
```

This separation keeps the engine portable and reusable.

---

# Philosophy

The engine should be:

- Deterministic
- Importable
- Framework agnostic
- Side-effect free
- Easy to test
- Easy to extract into its own repository in the future

The same inputs should always produce the same outputs.

---

# Concepts

## Facts

Facts are arbitrary pieces of data supplied by the caller.

Example:

```json
{
  "cart.total": 72,
  "customer.tier": "gold"
}
```

The engine does not understand what these values mean.

---

## Rules

Rules define:

- Conditions
- Actions

Example:

```json
{
  "id": "rule_10_percent_over_50",
  "conditions": {
    "all": [
      {
        "fact": "cart.total",
        "operator": "gte",
        "value": 50
      }
    ]
  },
  "actions": [
    {
      "type": "percentage_discount",
      "params": {
        "value": 10
      }
    }
  ]
}
```

---

## Actions

Actions are opaque to the engine.

Example:

```json
{
  "type": "percentage_discount",
  "params": {
    "value": 10
  }
}
```

The engine simply returns actions. The caller decides what they mean.

---

# Usage

```go
runtime := engine.New()

result, err := runtime.Evaluate(
    context.Background(),
    engine.EvaluationRequest{
        Facts: map[string]any{
            "cart.total": 72.0,
            "customer.tier": "gold",
        },
        Rules: rules,
    },
)
```

---

# Result

```go
type EvaluationResult struct {
    Actions        []Action
    MatchedRules   []string
    RejectedRules  []string
    Trace           []TraceEntry
}
```

---

# Trace

The engine provides a trace of its decisions.

Example:

```text
Rule: rule_10_percent_over_50
✓ cart.total >= 50
→ rule matched
```

This powers:

- Debugging
- Testing
- Simulation
- Explainability

---

# Supported Operators (V1)

- eq
- neq
- gt
- gte
- lt
- lte

---

# Future Capabilities

Potential future additions:

- Nested conditions
- OR groups
- NOT groups
- Arrays and collections
- String operators
- Custom operators
- Rule priorities
- Rule conflict resolution
- Rule compilation
- WASM support

---

# Non-Goals

The engine is not:

- A workflow engine
- A state machine
- A campaign manager
- A promotions platform

Those concerns belong to applications built on top of the engine.

---

# Guiding Principle

The engine answers one question:

> Given a set of facts and rules, which actions should be returned?


```go
runtime := engine.New()

result, _ := runtime.Evaluate(
    context.Background(),
    req,
)
```