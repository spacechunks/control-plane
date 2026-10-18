package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	jobendpointv1alpha1 "github.com/spacechunks/control-plane/internal/service/job/endpoint/v1alpha1"
	serverendpointv1alpha1 "github.com/spacechunks/control-plane/internal/service/server/endpoint/v1alpha1"
)

func main() {
	var (
		ctx           = context.Background()
		jobHandler    = jobendpointv1alpha1.NewHandler()
		serverHandler = serverendpointv1alpha1.NewHandler()
		logger        = slog.New(slog.NewJSONHandler(os.Stdout, nil))

		mux = http.NewServeMux()
	)

	if err := serverendpointv1alpha1.AddToMux(logger, serverHandler, mux); err != nil {
		die(logger, "failed adding server handlers", err)
	}

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
