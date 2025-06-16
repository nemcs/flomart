package repository

import (
	"context"
	"flomart/domain/user"
	"flomart/internal/identity"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO ⚠️	CreateUser возвращает *ID, а сервис возвращает ID. Выравнивай: указатели нужны только там, где nil — валидный вариант. Для id можно вернуть просто user.ID.
// TODO 🔧	FindByEmail: если пустой результат — верни pgx.ErrNoRows, а не fmt.Errorf(...), чтобы errors.Is работал.

type Repository interface {
	CreateUser(ctx context.Context, ser user.User) (*user.ID, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}
type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Создаем юзера в БД и возвращаем его id в service слой
// Работать с копией или указателем? Ничего не меняем, поэтому копией, а на практике в проде?
// ???Но чтобы возвращать nil, а не пустую структуру user.ID и нормально обработать пришлось добавить указатель, норм ли это практика?
func (repo *repository) CreateUser(ctx context.Context, u user.User) (*user.ID, error) {
	query := `
insert into users (email, password_hash, role, full_name, phone, is_active, created_at, updated_at) 
values ($1, $2, $3, $4, $5, $6, $7, $8) 
returning id`
	var id user.ID
	if err := repo.db.QueryRow(ctx, query, u.Email, u.PasswordHash, u.Role, u.FullName, u.Phone, u.IsActive, u.CreatedAt, u.UpdatedAt).Scan(&id); err != nil {
		return nil, fmt.Errorf("%s: %w", identity.ErrSqlInsertDev, err)
	}
	return &id, nil
}

func (repo *repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `select id, email, password_hash, role, full_name, phone, is_active, created_at, updated_at  from users where email = $1`
	row := repo.db.QueryRow(ctx, query, email)

	var u user.User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.FullName,
		&u.Phone,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", identity.ErrSqlSelectDev, err)
	}
	return &u, nil
}
