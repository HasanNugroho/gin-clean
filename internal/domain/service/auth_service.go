package service

import (
	"context"

	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (result dto.AuthResponse, err error)
	RefreshToken(ctx context.Context, request dto.RenewalTokenRequest) (result dto.AuthResponse, err error)
	Logout(ctx context.Context, accessToken string, request dto.RenewalTokenRequest) (err error)
}
