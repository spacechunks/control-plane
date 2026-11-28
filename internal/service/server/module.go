package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/spacechunks/control-plane/internal/queue"
	endpointv1alpha1 "github.com/spacechunks/control-plane/internal/service/server/endpoint/v1alpha1"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/postgres"
	"github.com/spacechunks/control-plane/internal/service/server/server"
	"github.com/spacechunks/control-plane/internal/service/server/worker"
)

type Config struct {
	DBURL            string          `env:"DB_URL,required"`
	Namespace        string          `env:"NAMESPACE,required"`
	Clusters         kube.ClusterSet `env:"CLUSTERS,required"`
	QueueConfig      queue.Config    `envPrefix:"QUEUE_CONFIG_"`
	MaxServerWorkers int             `env:"MAX_SERVER_WORKERS,required"`
}

var (
	q   *queue.RiverClient
	lis *server.PodDeletionListener
)

func Init(ctx context.Context, logger *slog.Logger, mux *http.ServeMux, cfg Config) error {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return fmt.Errorf("parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}

	var (
		pg = postgres.NewDB(logger.With("module", "server-postgres"), pool)
		qs = map[string]river.QueueConfig{
			"minigame_server_v1": {
				MaxWorkers: cfg.MaxServerWorkers,
			},
		}
	)

	workers := river.NewWorkers()

	createMinigame := worker.NewCreateMinigameServerWorker(
		logger.With("component", "create_minigame_server_worker"),
		pg,
		5*time.Minute, // TODO: configurable
		5*time.Minute, // TODO: configurable
		cfg.Clusters,
		cfg.Namespace,
	)

	if err := river.AddWorkerSafely(workers, createMinigame); err != nil {
		return fmt.Errorf("add create server worker: %w", err)
	}

	q, err = queue.NewRiverClient(logger.With("component", "river-server"), cfg.QueueConfig, pool, qs, workers)
	if err != nil {
		return fmt.Errorf("queue client: %w", err)
	}

	if err := q.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate riverqueue: %w", err)
	}

	lis := server.NewPodDeletionListener(
		logger.With("component", "pod-deletion-listener"),
		pg,
		cfg.Clusters,
		cfg.Namespace,
	)

	if err := lis.Start(ctx); err != nil {
		return fmt.Errorf("start pod deletion listener: %w", err)
	}

	var (
		serverService = server.NewService(
			logger.With("component", "server-service"),
			pg,
			pg,
			cfg.Clusters,
			cfg.Namespace,
			q,
		)
		serverHandler = endpointv1alpha1.NewHandler(serverService)
	)

	if err := endpointv1alpha1.AddToMux(logger, serverHandler, mux); err != nil {
		return fmt.Errorf("add handlers: %w", err)
	}

	logger.InfoContext(ctx, "migrating")

	if err := q.Start(ctx); err != nil {
		return fmt.Errorf("start queue: %w", err)
	}

	serverService.CreateMinigameServer(ctx, server.CreateMinigameServerArgs{
		Name:   "test",
		Labels: nil,
		Image:  "ghcr.io/spacechunks/paper-docker:26.2-112",
	})

	return nil
}

func Stop(ctx context.Context) error {
	lis.Stop()
	return q.Stop(ctx)
}
