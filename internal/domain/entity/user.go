package entity

import (
	"context"
	"time"

	"github.com/HasanNugroho/gin-clean/pkg/constants"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Email       string         `gorm:"uniqueIndex" json:"email"`
	PhoneNumber string         `gorm:"not null" json:"phone_number"`
	CipherText  string         `json:"-"`
	Role        constants.Role `gorm:"default:'user'" json:"role"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.CipherText), []byte(plainPassword))
	return err == nil
}

func (u *User) SetPassword(ctx context.Context, password string) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.CipherText = string(hash)
	return err
}
