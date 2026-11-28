package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/riverqueue/river/rivertype"
)

type InsertJobTXFunc func(
	ctx context.Context,
	tx pgx.Tx,
	args river.JobArgs,
	opts *river.InsertOpts,
) (*rivertype.JobInsertResult, error)

type Config struct {
	MaxAttempts int `env:"MAX_ATTEMPTS,required"`
}

type RiverClient struct {
	pool  *pgxpool.Pool
	river *river.Client[pgx.Tx]
}

func NewRiverClient(
	logger *slog.Logger,
	cfg Config,
	pool *pgxpool.Pool,
	queues map[string]river.QueueConfig,
	workers *river.Workers,
) (*RiverClient, error) {
	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues:      queues,
		Workers:     workers,
		MaxAttempts: cfg.MaxAttempts,
		Logger:      logger,
	})
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &RiverClient{
		pool:  pool,
		river: client,
	}, nil
}

func (c RiverClient) Start(ctx context.Context) error {
	return c.river.Start(ctx)
}

func (c RiverClient) Stop(ctx context.Context) error {
	return c.river.Stop(ctx)
}

func (c RiverClient) InsertJobTX(
	ctx context.Context,
	tx pgx.Tx,
	args river.JobArgs,
	opts *river.InsertOpts,
) (*rivertype.JobInsertResult, error) {
	return c.river.InsertTx(ctx, tx, args, opts)
}

func (c RiverClient) Migrate(ctx context.Context) error {
	lockTx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer lockTx.Rollback(context.Background())

	lockID := 67694201337

	// wait until any other instance running River migrations is finished.
	if _, err := lockTx.Exec(
		ctx,
		"SELECT pg_advisory_xact_lock($1)",
		lockID,
	); err != nil {
		return fmt.Errorf("lock: %w", err)
	}

	migrator, err := rivermigrate.New(riverpgxv5.New(c.pool), nil)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		return fmt.Errorf("up: %w", err)
	}

	if err := lockTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
