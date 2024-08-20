package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SqlxTransactionFn func(tx *sqlx.Tx) error

func Transaction(ctx context.Context, fn SqlxTransactionFn, db *DB) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("error rolling back transaction: %w", errRollback)
		}
		return fmt.Errorf("error executing a transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}
