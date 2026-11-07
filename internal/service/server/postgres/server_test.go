package postgres_test

import (
	"context"
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spacechunks/control-plane/internal/service/server/server"
	"github.com/spacechunks/control-plane/internal/service/server/test"
	"github.com/spacechunks/control-plane/internal/service/server/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestCreateServer(t *testing.T) {
	t.Parallel()

	var (
		ctx      = context.Background()
		pg       = runDB(t)
		expected = fixture.Server(func(s *server.Server) {
			s.Network = nil
		})
	)

	args := server.CreateServerArgs{
		ID:        expected.ID,
		Name:      expected.Name,
		Type:      expected.Type,
		State:     expected.State,
		Labels:    expected.Labels,
		CreatedAt: expected.CreatedAt,
		UpdatedAt: expected.UpdatedAt,
	}

	actual, err := pg.CreateServer(ctx, args)
	require.NoError(t, err)

	if d := cmp.Diff(expected, actual, cmpopts.EquateComparable(netip.Addr{})); d != "" {
		t.Fatalf("mismatch (-want +got):\n%s", d)
	}
}

func TestServerByID(t *testing.T) {
	t.Parallel()

	var (
		ctx      = context.Background()
		pg       = runDB(t)
		expected = fixture.Server()
	)

	args := server.CreateServerArgs{
		ID:        expected.ID,
		Name:      expected.Name,
		Type:      expected.Type,
		State:     expected.State,
		Labels:    expected.Labels,
		CreatedAt: expected.CreatedAt,
		UpdatedAt: expected.UpdatedAt,
	}

	test.InsertServer(t, expected)

	actual, err := pg.ServerByID(ctx, args.ID)
	require.NoError(t, err)

	if d := cmp.Diff(expected, actual, cmpopts.EquateComparable(netip.Addr{})); d != "" {
		t.Fatalf("mismatch (-want +got):\n%s", d)
	}
}
