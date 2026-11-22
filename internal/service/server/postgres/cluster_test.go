package postgres_test

import (
	"context"
	"testing"
	"uuid"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/server"
	"github.com/spacechunks/control-plane/internal/service/server/test"
	"github.com/spacechunks/control-plane/internal/service/server/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestBestClusterSelectsAndExhaustsSlots(t *testing.T) {
	var (
		ctx = context.Background()
		pg  = runDB(t)
		c1  = fixture.Cluster(func(c *kube.Cluster) {
			c.ID = uuid.New().String()
			c.Name = "c1"
			c.Slots = 1
		})
		c2 = fixture.Cluster(func(c *kube.Cluster) {
			c.ID = uuid.New().String()
			c.Name = "c2"
			c.Slots = 1
		})
	)

	test.InsertClusters(t, c1, c2)

	cluster1, err := pg.BestCluster(ctx)
	require.NoError(t, err)

	test.InsertServers(t, fixture.Server(func(s *server.Server) {
		s.ID = uuid.New().String()
		s.Name = "s1"
		s.Cluster = cluster1
	}))

	cluster2, err := pg.BestCluster(ctx)
	require.NoError(t, err)

	test.InsertServers(t, fixture.Server(func(s *server.Server) {
		s.ID = uuid.New().String()
		s.Name = "s2"
		s.Cluster = cluster2
	}))

	_, err = pg.BestCluster(ctx)
	require.ErrorIs(t, err, kube.ErrNoSlotsAvailable)
}
