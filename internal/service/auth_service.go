package service

import (
	"context"
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
)

type AuthService struct {
	repo           repository.UserRepository
	logger         *logger.Logger
	config         *config.Config
	jwt            *jwt.TokenGenerator
	contextTimeout time.Duration
}

func NewAuthService(repo repository.UserRepository, logger *logger.Logger, config *config.Config, jwt *jwt.TokenGenerator, timeout time.Duration) *AuthService {
	return &AuthService{
		repo:           repo,
		logger:         logger,
		config:         config,
		jwt:            jwt,
		contextTimeout: timeout,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (result dto.AuthResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.GetByEmail(ctx, req.Email)

	if err != nil {
		return result, errors.ErrUnauthorized.WithMessage("invalid email or password")
	}

	if !user.VerifyPassword(req.Password) {
		return result, errors.ErrUnauthorized.WithMessage("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return result, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return result, err
	}

	return dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Data: map[string]interface{}{
			"id": user.ID,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RenewalTokenRequest) (result dto.AuthResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	claims, err := s.jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return result, errors.ErrBadRequest.WithMessage("invalid refresh token").WithError(err)
	}

	if nbfRaw, ok := claims["nbf"]; ok {
		if nbf, ok := nbfRaw.(float64); ok && time.Now().Before(time.Unix(int64(nbf), 0)) {
			return result, errors.ErrUnauthorized.WithMessage("refresh token not yet valid")
		}
	}

	id, ok := claims["payload"].(string)
	if !ok {
		return result, errors.ErrUnauthorized.WithMessage("invalid token payload")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return result, errors.ErrBadRequest.WithMessage("user not found")
	}

	// Generate new access token
	newToken, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return result, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return result, err
	}

	result = dto.AuthResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		Data: map[string]interface{}{
			"id": user.ID,
		},
	}
	return
}

func (h *AuthService) Logout(c context.Context, accessToken string, req dto.RenewalTokenRequest) (err error) {
	// Blacklist Access Token
	if err = h.jwt.RevokeToken(accessToken); err != nil {
		return errors.ErrInternalServer.WithMessage("failed to revoke access token")
	}

	// Blacklist Refresh Token
	if err = h.jwt.RevokeRequestToken(req.RefreshToken); err != nil {
		return errors.ErrInternalServer.WithMessage("failed to revoke refresh token").WithError(err)
	}
	return nil
}
