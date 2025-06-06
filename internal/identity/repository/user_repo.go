package repository

import (
	"context"
	"flomart/domain/user"
	"flomart/internal/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ser identity.RegisterInput) (user.ID, error)
}
type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Создаем юзера в БД и возвращаем его id в service слой
func (repo *repository) CreateUser(u identity.RegisterInput) (user.ID, error) {
	query := `
insert into users (email, password_hash, full_name, phone, role) 
values ($1, $2, $3, $4, $5) 
returning id`
	var id user.ID
	if err := repo.db.QueryRow(context.Background(), query, u.Email, u.Password, u.Name, u.Phone, "client").Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
