package db

import (
	"context"
	"flomart/internal/identity"
	"flomart/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"time"
)

func NewPgxPool(dbURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Log.Error(identity.ErrDBConnectionDev,
			slog.String(logger.FieldErr, err.Error()),
			slog.String("dbURL", dbURL))
		os.Exit(1)
	}

	return pool
}
