package server

import (
	"context"
	"log/slog"

	"github.com/spacechunks/control-plane/internal/service/server/resource"
)

type CreateEphemeralServerArgs struct {
	Name   string
	Labels map[string]string
	// flavor
	// image
}

type Service interface {
	CreateEphemeralServer(ctx context.Context, args CreateEphemeralServerArgs) (resource.Server, error)
}

type Backend struct {
	logger *slog.Logger
}

func (s Backend) CreateEphemeralServer(ctx context.Context, args CreateEphemeralServerArgs) (resource.Server, error) {
	// TODO: db calls
	return resource.Server{}, nil
}

func NewService(logger *slog.Logger) *Backend {
	return &Backend{
		logger: logger,
	}
}
