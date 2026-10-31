package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
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
