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

type serverRelationsRow struct {
	Server  query.Server
	Cluster query.Cluster
}

func (db DB) CreateServer(ctx context.Context, args server.CreateServerDBArgs) (server.Server, error) {
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
		ClusterID: new(args.ClusterID),
	}

	if err := q.CreateServer(ctx, params); err != nil {
		return server.Server{}, fmt.Errorf("create: %w", err)
	}

	r, err := q.ServerByID(ctx, args.ID)
	if err != nil {
		return server.Server{}, fmt.Errorf("server by id: %w", err)
	}

	ret, err := toDomainServer(serverRelationsRow{
		Server:  r.Server,
		Cluster: r.Cluster,
	})
	if err != nil {
		return server.Server{}, fmt.Errorf("to domain: %w", err)
	}

	return ret, nil
}

func (db DB) ServerByID(ctx context.Context, id string) (server.Server, error) {
	release, q, err := db.queries(ctx)
	if err != nil {
		return server.Server{}, fmt.Errorf("queries: %w", err)
	}
	defer release()

	r, err := q.ServerByID(ctx, id)
	if err != nil {
		return server.Server{}, fmt.Errorf("get: %w", err)
	}

	ret, err := toDomainServer(serverRelationsRow{
		Server:  r.Server,
		Cluster: r.Cluster,
	})
	if err != nil {
		return server.Server{}, fmt.Errorf("to domain: %w", err)
	}

	return ret, nil
}

func (db DB) UpdateServer(ctx context.Context, args server.UpdateServerDBArgs) (server.Server, error) {
	var ret server.Server

	if err := db.doTX(ctx, func(q *query.Queries) error {
		old, err := q.ServerByID(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("server by id: %w", err)
		}

		update := query.UpdateServerParams{
			Name:      old.Server.Name,
			Labels:    old.Server.Labels,
			Ipv4:      old.Server.Ipv4,
			Port:      old.Server.Port,
			UpdatedAt: args.UpdatedAt,
			ID:        old.Server.ID,
			State:     old.Server.State,
		}

		if args.State != nil {
			update.State = query.ServerState(strings.ToUpper(args.State.String()))
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

		r, err := q.ServerByID(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("server by id after update: %w", err)
		}

		ret, err = toDomainServer(serverRelationsRow{
			Server:  r.Server,
			Cluster: r.Cluster,
		})
		if err != nil {
			return fmt.Errorf("to domain: %w", err)
		}

		return nil
	}); err != nil {
		return server.Server{}, err
	}

	return ret, nil
}

func (db DB) DeleteServer(ctx context.Context, id string) (server.Server, error) {
	var ret server.Server
	if err := db.doTX(ctx, func(q *query.Queries) error {
		r, err := q.ServerByID(ctx, id)
		if err != nil {
			return fmt.Errorf("server by id: %w", err)
		}

		if err := q.DeleteServer(ctx, id); err != nil {
			return fmt.Errorf("update: %w", err)
		}

		ret, err = toDomainServer(serverRelationsRow{
			Server:  r.Server,
			Cluster: r.Cluster,
		})
		if err != nil {
			return fmt.Errorf("to domain: %w", err)
		}

		return nil
	}); err != nil {
		return server.Server{}, err
	}

	return ret, nil
}
