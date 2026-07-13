package usecase_auth

import (
	"context"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/admin/auth/usecase/auth_login"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	AuthLoginUsecase interface {
		Execute(ctx context.Context, input domain.AuthLoginUsecaseInput) (*domain.AuthLoginUsecaseOutput, error)
	}
)

type (
	UserStore interface {
		Find(ctx context.Context, email string, orgID int64) (domain.User, error)
	}
	OrgStore interface {
		Find(ctx context.Context, publicID string) (domain.Organization, error)
	}
	PasswordSvc interface {
		Verify(password, hash string) bool
	}
	TokenSvc interface {
		Create(userID, orgID string) (string, time.Time, error)
	}
)

type AuthUsecaseFactory struct {
	authLoginUsecase AuthLoginUsecase
}

func NewUserUsecaseFactory(userStore UserStore, orgStore OrgStore, tokenSvc TokenSvc, passwordSvc PasswordSvc) *AuthUsecaseFactory {
	return &AuthUsecaseFactory{
		authLoginUsecase: auth_login.NewAuthLoginUsecase(userStore, orgStore, passwordSvc, tokenSvc),
	}
}

func (uc *AuthUsecaseFactory) LoginUserUsecase() AuthLoginUsecase {
	return uc.authLoginUsecase
}
