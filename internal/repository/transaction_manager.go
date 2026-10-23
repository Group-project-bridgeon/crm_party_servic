package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/group-project/crm_party_service/internal/db"
)

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(queries *db.Queries) error) error
}

type postgresTransactionManager struct {
	pool *pgxpool.Pool
}

func NewPostgresTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &postgresTransactionManager{pool: pool}
}

func (tm *postgresTransactionManager) ExecTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)
	if err := fn(q); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
