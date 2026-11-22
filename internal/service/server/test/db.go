package test

import (
	"context"
	"strings"
	"testing"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/server"
	testdb "github.com/spacechunks/control-plane/internal/test/db"
	"github.com/stretchr/testify/require"
)

func InsertServers(t *testing.T, servers ...server.Server) {
	var (
		ctx = context.Background()
		q   = `
INSERT INTO server
    (id, type, name, state, labels,created_at, updated_at, cluster_id)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8);
`
		conn = testdb.New(t)
	)

	for _, s := range servers {
		_, err := conn.Exec(
			ctx,
			q,
			s.ID,
			strings.ToUpper(s.Type.String()),
			s.Name,
			strings.ToUpper(s.State.String()),
			s.Labels,
			s.CreatedAt,
			s.UpdatedAt,
			s.Cluster.ID,
		)
		require.NoError(t, err)

		if s.Network != nil {
			_, err := conn.Exec(
				ctx,
				"UPDATE server SET ipv4 = $1, port = $2 WHERE id = $3",
				s.Network.IP4,
				s.Network.Port,
				s.ID,
			)
			require.NoError(t, err)
		}
	}
}

func InsertClusters(t *testing.T, clusters ...kube.Cluster) {
	var (
		ctx = context.Background()
		q   = `
			INSERT INTO cluster (id, name, slots, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
`
		conn = testdb.New(t)
	)

	for _, c := range clusters {
		_, err := conn.Exec(ctx, q, c.ID, c.Name, c.Slots, c.CreatedAt, c.UpdateAt)
		require.NoError(t, err)
	}
}
