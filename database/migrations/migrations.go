package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//TODO 🔧	Вместо fmt.Errorf("неверный source or url: %w", err) — верни AppError или заверши Fatal.

func RunMigrations(dbURL string) (*migrate.Migrate, error) {
	m, err := migrate.New("file://database/migrations", dbURL)
	if err != nil {
		return nil, fmt.Errorf("неверный source or url: %w\n", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("миграции не поднялись: %w\n", err)
	}
	return m, nil
}
