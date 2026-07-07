package api

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

func (s *handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	var httpReq EvaluateRequest

	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	engineReq := toEngineRequest(httpReq)

	result, err := s.engine.Evaluate(r.Context(), engineReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "evaluation_failed",
		})
		return
	}

	writeJSON(w, http.StatusOK, toEvaluateResponse(result))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
