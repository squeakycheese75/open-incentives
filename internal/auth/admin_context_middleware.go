package auth

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type (
	OrgContextStore interface {
		Find(rctx context.Context, publicID string) (domain.Organization, error)
	}
	ApiKeyContextStore interface {
		FindByPublicID(ctx context.Context, publicID string) (domain.APIKey, error)
	}
)

type cachedAuthContext struct {
	authCtx   AuthContext
	expiresAt time.Time
}

type AuthContextCache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]cachedAuthContext
}

func NewAuthContextCache(ttl time.Duration) *AuthContextCache {
	return &AuthContextCache{
		ttl:   ttl,
		items: make(map[string]cachedAuthContext),
	}
}

func (c *AuthContextCache) Get(key string) (AuthContext, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(item.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}
		return AuthContext{}, false
	}

	return item.authCtx, true
}

func (c *AuthContextCache) Set(key string, authCtx AuthContext) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cachedAuthContext{
		authCtx:   authCtx,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func AdminContextMiddleware(orgStore OrgContextStore, cache *AuthContextCache) func(http.Handler) http.Handler {
	if cache == nil {
		cache = NewAuthContextCache(5 * time.Minute)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, "missing auth claims")
				return
			}

			cacheKey := claims.OrgPublicID + ":" + claims.UserPublicID

			if authCtx, ok := cache.Get(cacheKey); ok {
				next.ServeHTTP(w, r.WithContext(ContextWithAuth(r.Context(), authCtx)))
				return
			}

			org, err := orgStore.Find(r.Context(), claims.OrgPublicID)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid organization")
				return
			}

			// user, err := store.FindUserByPublicID(r.Context(), claims.UserPublicID, org.ID)
			// if err != nil {
			// 	httputil.WriteError(w, http.StatusUnauthorized, "invalid user")
			// 	return
			// }

			authCtx := AuthContext{
				OrgID:       org.ID,
				OrgPublicID: org.PublicID,
				// UserID:       user.ID,
				// UserPublicID: user.PublicID,
			}

			cache.Set(cacheKey, authCtx)

			next.ServeHTTP(w, r.WithContext(ContextWithAuth(r.Context(), authCtx)))
		})
	}
}
