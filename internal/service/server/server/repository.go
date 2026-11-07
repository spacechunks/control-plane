package server

import (
	"context"
	"time"
)

type CreateServerArgs struct {
	ID        string
	Name      string
	Labels    map[string]string
	Type      ServerType
	State     ServerState
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateServerArgs struct {
	ID        string
	Name      *string
	Labels    map[string]string
	Network   *Network
	UpdatedAt time.Time
}

type Repository interface {
	CreateServer(ctx context.Context, args CreateServerArgs) error
	ServerByID(ctx context.Context, id string) (Server, error)
	UpdateServer(ctx context.Context, args UpdateServerArgs) (Server, error)
}
