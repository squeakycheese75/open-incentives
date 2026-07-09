package auth_login

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/domain"
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

type AuthLoginUsecase struct {
	userStore   UserStore
	orgStore    OrgStore
	passwordSvc PasswordSvc
	tokenSvc    TokenSvc
}

func NewAuthLoginUsecase(userStore UserStore, orgStore OrgStore, passwordSvc PasswordSvc, tokenSvc TokenSvc) *AuthLoginUsecase {
	return &AuthLoginUsecase{
		userStore:   userStore,
		orgStore:    orgStore,
		passwordSvc: passwordSvc,
		tokenSvc:    tokenSvc,
	}
}

func (uc *AuthLoginUsecase) Execute(ctx context.Context, input domain.AuthLoginUsecaseInput) (*domain.AuthLoginUsecaseOutput, error) {
	orgSlug := strings.TrimSpace(input.OrgPublicID)
	if orgSlug == "" {
		return nil, errors.New("org slug is missing")
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	org, err := uc.orgStore.Find(ctx, orgSlug)
	if err != nil {
		return nil, err
	}

	user, err := uc.userStore.Find(ctx, email, org.ID)
	if err != nil {
		return nil, err
	}

	if !uc.passwordSvc.Verify(input.Password, user.HashedPassword) {
		return nil, errors.New("invalid credentials")
	}

	// // Generate a new JWT Token
	token, expiry, err := uc.tokenSvc.Create(user.PublicID, orgSlug)
	if err != nil {
		return nil, errors.New("token generation failed")
	}

	return &domain.AuthLoginUsecaseOutput{Token: token, ExpiresAt: expiry}, err
}
