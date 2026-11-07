package server

import (
	"context"
	"log/slog"
)

type CreateEphemeralServerArgs struct {
	Name   string
	Labels map[string]string
	// flavor
	// image
}

type Service interface {
	CreateEphemeralServer(ctx context.Context, args CreateEphemeralServerArgs) (Server, error)
}

type Backend struct {
	logger *slog.Logger
}

func (s Backend) CreateEphemeralServer(ctx context.Context, args CreateEphemeralServerArgs) (Server, error) {
	// TODO: db calls

	//fixture.Server()
	return Server{}, nil
}

func NewService(logger *slog.Logger) *Backend {
	return &Backend{
		logger: logger,
	}
}
