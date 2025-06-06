// storage - Взаимодействие с БД (pgx)
// Работа с таблицей users
package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func NewPool(ctx context.Context, conn string) {
	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		log.Fatal("Unable to create connection pool: " + err.Error())
	}
	defer dbpool.Close()
}
