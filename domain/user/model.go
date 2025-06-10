package user

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ID string

type User struct {
	ID           ID        `json:"ID"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
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
