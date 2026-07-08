package eval

import (
	"context"

	engine "github.com/squeakycheese75/open-incentives-engine"
)

type (
	Engine interface {
		Evaluate(ctx context.Context, req engine.EvaluationRequest) (engine.EvaluationResult, error)
	}
)

type Handler struct {
	engine Engine
}

func NewHandler(engine Engine) *Handler {
	return &Handler{
		engine: engine,
	}
}
