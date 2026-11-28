package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	jobendpointv1alpha1 "github.com/spacechunks/control-plane/internal/service/job/endpoint/v1alpha1"
	"github.com/spacechunks/control-plane/internal/service/server"
)

type config struct {
	ServerConfig server.Config `envPrefix:"SERVER_"`
}

func main() {
	var (
		ctx    = context.Background()
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		mux    = http.NewServeMux()

		// TODO: move to job.Init
		jobHandler = jobendpointv1alpha1.NewHandler()
	)

	cfg, err := env.ParseAs[config]()
	if err != nil {
		die(logger, "failed to read config", err)
	}

	// service inits

	if err := server.Init(ctx, logger, mux, cfg.ServerConfig); err != nil {
		die(logger, "failed to initialize server service", err)
	}

	// TODO: move to job.Init
	if err := jobendpointv1alpha1.AddToMux(logger, jobHandler, mux); err != nil {
		die(logger, "failed adding job handlers", err)
	}

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		_ = serv.Shutdown(ctx)

		if err := server.Stop(ctx); err != nil {
			logger.ErrorContext(ctx, "error stopping server module", "err", err)
		}
	}()

	if err := serv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			die(logger, "failed to stop server", err)
		}
	}
}

func die(logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "err", err)
	os.Exit(1)
}
