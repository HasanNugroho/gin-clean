package service

import (
	"context"

	"github.com/HasanNugroho/gin-clean/internal/domain/entity"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
)

type UserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (err error)
	GetById(ctx context.Context, id string) (user *entity.User, err error)
	GetByEmail(ctx context.Context, email string) (user *entity.User, err error)
	Update(ctx context.Context, id string, user *dto.UpdateUserRequest) (err error)
	Delete(ctx context.Context, id string) (err error)
}
