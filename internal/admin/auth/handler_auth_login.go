package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type AuthLoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

type AuthLoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (r *AuthLoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("missing_email")
	}

	if r.Password == "" {
		return errors.New("missing_password")
	}
	if r.Organization == "" {
		return errors.New("missing_organization")
	}

	return nil
}

func (s *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	if err := req.Validate(); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	res, err := s.authContainer.LoginUserUsecase().Execute(r.Context(), domain.AuthLoginUsecaseInput{
		Email:       req.Email,
		Password:    req.Password,
		OrgPublicID: req.Organization,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrUnauthorized):
			httputil.WriteJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_login_user"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, AuthLoginResponse{
		Token:     res.Token,
		ExpiresAt: res.ExpiresAt,
	})
}
