package server

import "context"

type CreateServerArgs struct {
	ID     string
	Name   string
	Labels map[string]string
	Type   string
}

type Repository interface {
	CreateServer(ctx context.Context, args CreateServerArgs) error
}
