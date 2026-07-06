package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	engine "github.com/squeakycheese75/open-incentives-engine"
)

var ErrInvalidRequest = errors.New("invalid request")

type Engine interface {
	Evaluate(ctx context.Context, req engine.EvaluationRequest) (engine.EvaluationResult, error)
}

type handlers struct {
	engine Engine
}

func NewHandlers(engine Engine) *handlers {
	return &handlers{
		engine: engine,
	}
}

func (s *handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
