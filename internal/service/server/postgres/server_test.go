package postgres_test

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5"
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

func TestUpdateServer(t *testing.T) {
	now := time.Now()
	tc := []struct {
		name     string
		args     server.UpdateServerArgs
		expected server.Server
	}{
		{
			name: "update name",
			args: server.UpdateServerArgs{
				ID:        fixture.Server().ID,
				Name:      new("new-name"),
				UpdatedAt: now,
			},
			expected: fixture.Server(func(s *server.Server) {
				s.Name = "new-name"
				s.UpdatedAt = now
			}),
		},
		{
			name: "update labels",
			args: server.UpdateServerArgs{
				ID:        fixture.Server().ID,
				Labels:    map[string]string{"new": "labels"},
				UpdatedAt: now,
			},
			expected: fixture.Server(func(s *server.Server) {
				s.Labels = map[string]string{"new": "labels"}
				s.UpdatedAt = now
			}),
		},
		{
			name: "update network",
			args: server.UpdateServerArgs{
				ID: fixture.Server().ID,
				Network: &server.Network{
					IP4:  netip.MustParseAddr("192.0.2.33"),
					Port: 1337,
				},
				UpdatedAt: now,
			},
			expected: fixture.Server(func(s *server.Server) {
				s.Network = &server.Network{
					IP4:  netip.MustParseAddr("192.0.2.33"),
					Port: 1337,
				}
				s.UpdatedAt = now
			}),
		},
		{
			name: "update all",
			args: server.UpdateServerArgs{
				ID:     fixture.Server().ID,
				Name:   new("new-name"),
				Labels: map[string]string{"2": "3"},
				Network: &server.Network{
					IP4:  netip.MustParseAddr("192.0.2.33"),
					Port: 1337,
				},
				UpdatedAt: now,
			},
			expected: fixture.Server(func(s *server.Server) {
				s.Name = "new-name"
				s.Labels = map[string]string{"2": "3"}
				s.Network = &server.Network{
					IP4:  netip.MustParseAddr("192.0.2.33"),
					Port: 1337,
				}
				s.UpdatedAt = now
			}),
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				ctx = context.Background()
				pg  = runDB(t)
			)

			test.InsertServer(t, fixture.Server())

			actual, err := pg.UpdateServer(ctx, tt.args)
			require.NoError(t, err)

			if d := cmp.Diff(tt.expected, actual, cmpopts.EquateComparable(netip.Addr{})); d != "" {
				t.Fatalf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestDeleteServer(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.Background()
		pg       = runDB(t)
		expected = fixture.Server()
	)

	test.InsertServer(t, expected)

	actual, err := pg.DeleteServer(ctx, expected.ID)
	require.NoError(t, err)

	if d := cmp.Diff(expected, actual, cmpopts.EquateComparable(netip.Addr{})); d != "" {
		t.Fatalf("mismatch (-want +got):\n%s", d)
	}

	_, err = pg.ServerByID(ctx, expected.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
