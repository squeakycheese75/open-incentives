package auth

import usecase_auth "github.com/squeakycheese75/open-incentives/internal/admin/auth/usecase"

type Handler struct {
	authContainer *usecase_auth.AuthUsecaseFactory
}

func NewHandler(authContainer *usecase_auth.AuthUsecaseFactory) *Handler {
	return &Handler{
		authContainer: authContainer,
	}
}
