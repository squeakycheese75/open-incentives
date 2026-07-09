package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const adminClaimsKey contextKey = "admin_claims"

func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, adminClaimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(adminClaimsKey).(Claims)
	return claims, ok
}

type Claims struct {
	UserID string `json:"user_id"`
	OrgID  string `json:"org_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}
