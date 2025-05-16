package postgresql

import (
	"context"

	"github.com/HasanNugroho/gin-clean/internal/domain/entity"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	db := u.db.WithContext(ctx)

	result := db.Create(user)
	if result.Error != nil {
		return errors.Wrap(errors.ErrInternalServer, result.Error)
	}

	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	db := u.db.WithContext(ctx)

	var user entity.User
	result := db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrInternalServer, result.Error)
	}

	return &user, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := u.db.WithContext(ctx)

	var user entity.User
	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrInternalServer, result.Error)
	}

	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *entity.User) error {
	db := u.db.WithContext(ctx)

	result := db.Save(user)
	if result.Error != nil {
		return errors.Wrap(errors.ErrInternalServer, result.Error)
	}

	return nil
}

func (u *userRepository) Delete(ctx context.Context, id string) error {
	db := u.db.WithContext(ctx)

	result := db.Delete(&entity.User{}, id)
	if result.Error != nil {
		return errors.Wrap(errors.ErrInternalServer, result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}
