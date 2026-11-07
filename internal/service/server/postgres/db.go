package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spacechunks/control-plane/internal/service/server/postgres/query"
)

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
