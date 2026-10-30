package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const tokenLocal = "token"

var ErrInvalidToken = errors.New("invalid or expired token")

func Issue(secret []byte, subject string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func Verify(secret []byte, raw string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(*jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func Middleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		raw, found := strings.CutPrefix(header, "Bearer ")
		if !found || raw == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}
		if _, err := Verify(secret, raw); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		c.Locals(tokenLocal, raw)
		return c.Next()
	}
}

func Token(c *fiber.Ctx) string {
	raw, _ := c.Locals(tokenLocal).(string)
	return raw
}
