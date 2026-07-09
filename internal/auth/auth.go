package auth

import "context"

type AuthContext struct {
	OrgID        int64
	OrgPublicID  string
	UserID       int64
	UserPublicID string
}

type authContextKey struct{}

func ContextWithAuth(ctx context.Context, authCtx AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey{}, authCtx)
}

func AuthFromContext(ctx context.Context) (AuthContext, bool) {
	authCtx, ok := ctx.Value(authContextKey{}).(AuthContext)
	return authCtx, ok
}
