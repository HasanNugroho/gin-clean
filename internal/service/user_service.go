package service

import (
	"context"
	"time"

	"github.com/HasanNugroho/gin-clean/internal/domain/entity"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
)

type UserService struct {
	repo           repository.UserRepository
	contextTimeout time.Duration
}

func NewUserService(repo repository.UserRepository, timeout time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return errors.Wrap(errors.ErrConflict, err)
	}

	user := &entity.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
		IsActive:    true,
	}

	if err := user.SetPassword(ctx, req.Password); err != nil {
		return errors.Wrap(errors.ErrBadRequest, err)
	}

	return u.repo.Create(ctx, user)
}

func (u *UserService) GetById(ctx context.Context, id string) (user *entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByEmail(ctx, email)
}

func (u *UserService) Update(ctx context.Context, id string, updatedUser *dto.UpdateUserRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(errors.ErrNotFound, err)
	}
	if existing == nil {
		return errors.ErrNotFound
	}

	existing.Name = updatedUser.Name
	existing.Email = updatedUser.Email
	existing.PhoneNumber = updatedUser.PhoneNumber
	existing.Role = updatedUser.Role
	existing.IsActive = updatedUser.IsActive
	existing.UpdatedAt = time.Now()

	return u.repo.Update(ctx, existing)
}

func (u *UserService) Delete(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(errors.ErrNotFound, err)
	}
	if existing == nil {
		return errors.ErrNotFound
	}

	return u.repo.Delete(ctx, id)
}
