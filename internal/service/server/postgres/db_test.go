package postgres_test

import (
	"testing"

	"github.com/spacechunks/control-plane/internal/service/server/postgres"
	"github.com/spacechunks/control-plane/internal/service/server/postgres/migrations"
	testdb "github.com/spacechunks/control-plane/internal/test/db"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testdb.Run()
	m.Run()
	testdb.Run()
}

func runDB(t *testing.T) *postgres.DB {
	conn := testdb.New(t)
	err := migrations.Migrate(conn.Config().ConnString())
	require.NoError(t, err)
	return postgres.NewDB(nil, conn)
}
