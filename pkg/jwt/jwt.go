package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

type (
	TokenGenerator struct {
		cache               *cache.RedisCache
		secret              []byte
		tokenExpired        time.Duration
		refreshTokenExpired time.Duration
	}
)

func SetJWTHelper(config *config.Config, redis *cache.RedisCache) *TokenGenerator {
	tokenExpiry, _ := time.ParseDuration(config.Secret.TokenExpiry)
	refreshExpiry, _ := time.ParseDuration(config.Secret.RefreshTokenExpiry)
	return &TokenGenerator{
		cache:               redis,
		secret:              []byte(config.Secret.Jwt),
		tokenExpired:        tokenExpiry,
		refreshTokenExpired: refreshExpiry,
	}
}

func (t *TokenGenerator) GenerateToken(payload string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(t.tokenExpired).Unix(),
		"iat":     time.Now().Unix(),
	})

	return claims.SignedString(t.secret)
}

func (t *TokenGenerator) GenerateRefreshToken(payload string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(t.refreshTokenExpired).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Add(t.tokenExpired).Unix(),
	})

	return claims.SignedString(t.secret)
}

func (t *TokenGenerator) ParseToken(rawToken string) (jwt.MapClaims, error) {
	if t.IsTokenRevoked("refreshtoken:blacklist:"+rawToken) || t.IsTokenRevoked("token:blacklist:"+rawToken) {
		return nil, errors.ErrUnauthorized.WithMessage("invalid or expired token")
	}

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrUnauthorized.WithMessage("unexpected signing method")
		}
		return t.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.ErrUnauthorized.WithMessage("invalid or expired token").WithError(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && token.Valid {
		return nil, errors.ErrUnauthorized.WithMessage("invalid token claims")
	}
	return claims, nil
}

func (t *TokenGenerator) ParseRefreshToken(rawToken string) (jwt.MapClaims, error) {
	if t.IsTokenRevoked("refreshtoken:blacklist:"+rawToken) || t.IsTokenRevoked("token:blacklist:"+rawToken) {
		return nil, errors.ErrUnauthorized.WithMessage("invalid or expired token")
	}

	token, _, err := jwt.NewParser().ParseUnverified(rawToken, jwt.MapClaims{})
	if err != nil {
		return nil, errors.ErrUnauthorized.WithMessage("failed to parse refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && token.Valid {
		return nil, errors.ErrUnauthorized.WithMessage("invalid token claims")
	}
	return claims, nil
}

func (t *TokenGenerator) RevokeToken(token string) error {
	claims, err := t.ParseToken(token)
	if err != nil {
		return err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.ErrUnauthorized.WithMessage("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = t.cache.Set(context.Background(), "token:blacklist:"+token, "revoked", ttl); err != nil {
		return errors.ErrUnauthorized.WithMessage("failed to store token in blacklist")
	}
	fmt.Println("Token blacklisted:", token, "TTL:", ttl)

	return nil
}

func (t *TokenGenerator) RevokeRequestToken(token string) error {
	claims, err := t.ParseRefreshToken(token)
	if err != nil {
		return err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.ErrUnauthorized.WithMessage("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = t.cache.Set(context.Background(), "refreshtoken:blacklist:"+token, "revoked", ttl); err != nil {
		return errors.ErrUnauthorized.WithMessage("failed to store refresh token in blacklist")
	}
	fmt.Println("Token blacklisted:", token, "TTL:", ttl)

	return nil
}

func (t *TokenGenerator) IsTokenRevoked(tokenKey string) bool {
	val, err := t.cache.Exist(context.Background(), tokenKey)
	return err == nil && val > 0
}
