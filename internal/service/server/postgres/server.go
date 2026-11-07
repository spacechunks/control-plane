package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spacechunks/control-plane/internal/service/server/postgres/query"
	"github.com/spacechunks/control-plane/internal/service/server/server"
)

func (db DB) CreateServer(ctx context.Context, args server.CreateServerArgs) (server.Server, error) {
	release, q, err := db.queries(ctx)
	if err != nil {
		return server.Server{}, fmt.Errorf("queries: %w", err)
	}
	defer release()

	var labels []byte
	labels, err = json.Marshal(args.Labels)
	if err != nil {
		return server.Server{}, fmt.Errorf("marshal labels: %w", err)
	}

	params := query.CreateServerParams{
		ID:        args.ID,
		Type:      query.ServerType(strings.ToUpper(args.Type.String())),
		Name:      args.Name,
		State:     query.ServerState(strings.ToUpper(args.State.String())),
		Labels:    labels,
		CreatedAt: args.CreatedAt,
		UpdatedAt: args.UpdatedAt,
	}

	if err := q.CreateServer(ctx, params); err != nil {
		return server.Server{}, fmt.Errorf("create: %w", err)
	}

	return server.Server{
		ID:        args.ID,
		Name:      args.Name,
		Type:      args.Type,
		Labels:    args.Labels,
		State:     args.State,
		CreatedAt: args.CreatedAt,
		UpdatedAt: args.UpdatedAt,
	}, nil
}

func (db DB) ServerByID(ctx context.Context, id string) (server.Server, error) {
	release, q, err := db.queries(ctx)
	if err != nil {
		return server.Server{}, fmt.Errorf("queries: %w", err)
	}
	defer release()

	s, err := q.ServerByID(ctx, id)
	if err != nil {
		return server.Server{}, fmt.Errorf("get: %w", err)
	}

	var labels map[string]string
	if err := json.Unmarshal(s.Labels, &labels); err != nil {
		return server.Server{}, fmt.Errorf("unmarshal labels: %w", err)
	}

	ret := server.Server{
		ID:        s.ID,
		Name:      s.Name,
		Type:      server.ServerType(strings.ToLower(string(s.Type))),
		Labels:    labels,
		State:     server.ServerState(strings.ToLower(string(s.State))),
		CreatedAt: s.CreatedAt.UTC(),
		UpdatedAt: s.UpdatedAt.UTC(),
	}

	if s.Ipv4 != nil {
		ret.Network = &server.Network{
			IP4:  *s.Ipv4,
			Port: uint16(s.Port.Int32),
		}
	}

	return ret, nil
}
