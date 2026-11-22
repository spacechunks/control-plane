package server

import (
	"context"
	"time"
)

type CreateServerDBArgs struct {
	ID        string
	Name      string
	Labels    map[string]string
	Type      ServerType
	State     ServerState
	ClusterID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateServerDBArgs struct {
	ID        string
	Name      *string
	Labels    map[string]string
	Network   *Network
	State     *ServerState
	UpdatedAt time.Time
}

type Repository interface {
	CreateServer(ctx context.Context, args CreateServerDBArgs) (Server, error)
	ServerByID(ctx context.Context, id string) (Server, error)
	UpdateServer(ctx context.Context, args UpdateServerDBArgs) (Server, error)
	DeleteServer(ctx context.Context, id string) (Server, error)
}
