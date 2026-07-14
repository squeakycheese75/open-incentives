package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/squeakycheese75/open-incentives/internal/auth"
)

type TokenService struct {
	secret []byte
}

func NewJWTTokenService(secret string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
	}
}

func (s *TokenService) Create(userID, orgID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userPublicId": userID,
		"orgPublicId":  orgID,
		"role":         "admin",
		"exp":          expiresAt.Unix(),
	})

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (s *TokenService) Verify(tokenString string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *TokenService) HashToken(rawToken []byte) (string, error) {
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	return tokenHash, nil
}
