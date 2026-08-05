package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"maxwin/models"
)

const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
)

var (
	ErrMissingSecret = errors.New("JWT_SECRET is not set")
	ErrInvalidToken  = errors.New("invalid or expired token")
)

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func secret() ([]byte, error) {
	s := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if s == "" {
		return nil, ErrMissingSecret
	}
	return []byte(s), nil
}

func expiry() time.Duration {
	hours := 24
	if raw := strings.TrimSpace(os.Getenv("JWT_EXPIRY_HOURS")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			hours = n
		}
	}
	return time.Duration(hours) * time.Hour
}

// MintToken issues a short-lived JWT bound to the given user.
func MintToken(user models.User) (token string, expiresAt time.Time, err error) {
	key, err := secret()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt = time.Now().Add(expiry())
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(key)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}
	return signed, expiresAt, nil
}

// ParseToken validates a JWT and returns its claims.
func ParseToken(tokenString string) (*Claims, error) {
	key, err := secret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == "" || claims.Username == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// BearerToken extracts the raw token from an Authorization header.
func BearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", ErrInvalidToken
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", ErrInvalidToken
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", ErrInvalidToken
	}
	return token, nil
}
