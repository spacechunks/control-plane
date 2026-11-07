package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spacechunks/control-plane/internal/service/server/postgres/query"
	"github.com/spacechunks/control-plane/internal/service/server/server"
)

var _ server.Repository = (*DB)(nil)

type DB struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func NewDB(logger *slog.Logger, pool *pgxpool.Pool) *DB {
	return &DB{
		logger: logger,
		pool:   pool,
	}
}

func (db DB) queries(ctx context.Context) (func(), *query.Queries, error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("acquire conn: %w", err)
	}

	return conn.Release, query.New(conn), nil
}

func (db DB) doTX(ctx context.Context, fn func(q *query.Queries) error) error {
	if err := db.pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		if err := fn(query.New(tx)); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				// only log rollback related error, because we want to
				// pass the actual error that caused the rollback back to the caller
				db.logger.ErrorContext(ctx, "failed to rollback tx", "err", txErr)
			}
			return err
		}

		return tx.Commit(ctx)
	}); err != nil {
		return err
	}

	return nil
}
