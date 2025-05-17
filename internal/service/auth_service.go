package service

import (
	"context"
	"fmt"
	"net/http"
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
	contextTimeout time.Duration
}

func NewAuthService(repo repository.UserRepository, logger *logger.Logger, config *config.Config, timeout time.Duration) *AuthService {
	return &AuthService{
		repo:           repo,
		logger:         logger,
		config:         config,
		contextTimeout: timeout,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.GetByEmail(ctx, req.Email)
	fmt.Println(err)

	if err != nil {
		return dto.AuthResponse{}, errors.New("UNAUTHORIZED", "invalid email or password", http.StatusUnauthorized, err)
	}

	if !user.VerifyPassword(req.Password) {
		return dto.AuthResponse{}, errors.New("UNAUTHORIZED", "invalid email or password", http.StatusUnauthorized, err)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Data: map[string]interface{}{
			"id": user.ID,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RenewalTokenRequest) (dto.AuthResponse, error) {
	// ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	// defer cancel()
	// userID, err := jwt.ParseRefreshToken(req.RefreshToken, s.config.JWT.Secret)
	// if err != nil {
	// 	return dto.AuthResponse{}, errors.New("invalid refresh token")
	// }

	// user, err := s.repo.FindByID(ctx, userID)
	// if err != nil {
	// 	return dto.AuthResponse{}, errors.New("user not found")
	// }

	// // Generate new access token
	// newToken, err := jwtutil.GenerateToken(user.ID, s.config.JWT.Secret, time.Hour)
	// if err != nil {
	// 	return dto.AuthResponse{}, err
	// }

	// newRefreshToken, err := jwtutil.GenerateRefreshToken(user.ID, s.config.JWT.Secret, s.config.JWT.RefreshTokenExpiry)
	// if err != nil {
	// 	return dto.AuthResponse{}, err
	// }

	// return dto.AuthResponse{
	// 	Token:        newToken,
	// 	RefreshToken: newRefreshToken,
	// 	Data: map[string]interface{}{
	// 		"id":    user.ID,
	// 		"email": user.Email,
	// 		"name":  user.Name,
	// 	},
	// }, nil
	panic("un")
}
