package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func SafeRollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		fmt.Errorf("rollback error: %w", err)
	}
}
