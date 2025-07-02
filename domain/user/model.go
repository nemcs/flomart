package user

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ID string

type User struct {
	ID           ID        `json:"id" validate:"required,uuid4"`                        // UUID обязательно
	Email        string    `json:"email" validate:"required,email"`                     // Email обязателен и должен быть валидным
	PasswordHash string    `json:"-" validate:"required"`                               // Пароль должен быть хэширован, но обязателен
	Role         string    `json:"role" validate:"required,oneof=admin client courier"` // Роль обязательна, можно задать допустимые значения
	FullName     string    `json:"full_name" validate:"required,min=2,max=100"`         // Обязательное ФИО, с ограничением длины
	Phone        string    `json:"phone" validate:"required,e164"`                      // Телефон обязателен, e164 — формат вроде +71234567890
	IsActive     bool      `json:"is_active"`                                           // Логическое значение — не валидируется
	CreatedAt    time.Time `json:"created_at"`                                          // Эти два поля обычно не валидируются
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, c *User) error
}

func New(email, password, role, name, phone string) *User {
	return &User{
		ID:           ID(uuid.New().String()),
		Email:        email,
		PasswordHash: password,
		Role:         role,
		FullName:     name,
		Phone:        phone,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
