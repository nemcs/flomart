package db

import (
	"context"
	"flomart/pkg/logger"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

func SafeRollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		logger.Log.Error("rollback failed", slog.String(logger.FieldErr, err.Error()))
	}
}
