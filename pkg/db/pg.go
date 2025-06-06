package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func NewPgxPool(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	return pool
}
