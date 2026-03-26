// Package crypto provides JWT generation and parsing utilities using HS256.
package crypto

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skillora/backend/internal/config"
)

const jwtTTL = 7 * 24 * time.Hour // 7 days

// Claims represents the custom JWT payload.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed HS256 JWT containing user_id and exp.
func GenerateJWT(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTTL)),
			Issuer:    "skillora",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.C.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}
	return signed, nil
}

// ParseJWT validates a token string and returns the claims on success.
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("jwt: unexpected signing method %v", t.Header["alg"])
			}
			return []byte(config.C.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("jwt parse: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("jwt: invalid token")
	}
	if claims.UserID == "" {
		return nil, errors.New("jwt: missing user_id claim")
	}

	return claims, nil
}
