package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
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

	ret, err := toDomainServer(s)
	if err != nil {
		return server.Server{}, fmt.Errorf("to domain: %w", err)
	}

	return ret, nil
}

func (db DB) UpdateServer(ctx context.Context, args server.UpdateServerArgs) (server.Server, error) {
	var ret server.Server

	if err := db.doTX(ctx, func(q *query.Queries) error {
		old, err := q.ServerByID(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("server by id: %w", err)
		}

		update := query.UpdateServerParams{
			Name:      old.Name,
			Labels:    old.Labels,
			Ipv4:      old.Ipv4,
			Port:      old.Port,
			UpdatedAt: args.UpdatedAt,
			ID:        old.ID,
		}

		if args.Name != nil {
			update.Name = *args.Name
		}

		if args.Network != nil {
			update.Ipv4 = &args.Network.IP4
			update.Port = pgtype.Int4{
				Int32: int32(args.Network.Port),
				Valid: true,
			}
		}

		if args.Labels != nil {
			labels, err := json.Marshal(args.Labels)
			if err != nil {
				return fmt.Errorf("marshal labels: %w", err)
			}
			update.Labels = labels
		}

		if err := q.UpdateServer(ctx, update); err != nil {
			return fmt.Errorf("update: %w", err)
		}

		s, err := q.ServerByID(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("server by id after update: %w", err)
		}

		ret, err = toDomainServer(s)
		if err != nil {
			return fmt.Errorf("to domain: %w", err)
		}

		return nil
	}); err != nil {
		return server.Server{}, err
	}

	return ret, nil
}
