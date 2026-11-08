package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	SecretKey []byte
	Issuer    string
	Audience  string
	TTL       time.Duration
}

func GenerateAccessToken(cfg TokenConfig, userID string, role string) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID,

			IssuedAt: jwt.NewNumericDate(now),

			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.TTL)),

			Issuer: cfg.Issuer,

			Audience: jwt.ClaimStrings{cfg.Audience},

			ID: generateTokenID(),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(cfg.SecretKey)

	if err != nil {
		return "", fmt.Errorf("token could not be signed: %w", err)
	}

	return signedToken, nil
}

func generateTokenID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func ParseAccessToken(cfg TokenConfig, tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected algorithm: %s", t.Method.Alg())
			}
			return cfg.SecretKey, nil
		},
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("token time has expired: %w", err)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, fmt.Errorf("token format is incorrect: %w", err)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, fmt.Errorf("token signature is invalid: %w", err)
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, fmt.Errorf("token is not valid yet: %w", err)
		case errors.Is(err, jwt.ErrTokenInvalidIssuer):
			return nil, fmt.Errorf("token issuer is invalid: %w", err)
		case errors.Is(err, jwt.ErrTokenInvalidAudience):
			return nil, fmt.Errorf("token audience is invalid: %w", err)
		default:
			return nil, fmt.Errorf("token could not be parsed: %w", err)
		}
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
