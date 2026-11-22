package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	endpointv1alpha1 "github.com/spacechunks/control-plane/internal/service/server/endpoint/v1alpha1"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/postgres"
	"github.com/spacechunks/control-plane/internal/service/server/server"
)

type Config struct {
	DBURL     string          `env:"DB_URL,required"`
	Namespace string          `env:"NAMESPACE,required"`
	Clusters  kube.ClusterSet `env:"CLUSTERS,required"`
}

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
		pg            = postgres.NewDB(logger.With("module", "server-postgres"), pool)
		serverService = server.NewService(
			logger.With("component", "server-service"),
			pg,
			pg,
			cfg.Clusters,
			cfg.Namespace,
		)
		serverHandler = endpointv1alpha1.NewHandler(serverService)
	)

	if err := endpointv1alpha1.AddToMux(logger, serverHandler, mux); err != nil {
		return fmt.Errorf("add handlers: %w", err)
	}

	return nil
}
