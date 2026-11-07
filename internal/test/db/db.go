package db

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pg *handle
	mu sync.Mutex
)

type handle struct {
	dns string
	ctr *postgres.PostgresContainer
}

func Run() {
	mu.Lock()
	defer mu.Unlock()
	if pg != nil {
		return
	}

	ctx := context.Background()
	ctr, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	pg = &handle{
		dns: dsn,
		ctr: ctr,
	}
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if pg == nil {
		return
	}

	err := testcontainers.TerminateContainer(pg.ctr)
	if err != nil {
		panic(err)
	}

	pg = nil
}

func New(t *testing.T) *pgxpool.Pool {
	//mu.Lock()
	//defer mu.Unlock()
	if pg == nil {
		t.Fatal("no postgres container")
	}

	ctx := context.Background()
	admin, err := pgx.Connect(ctx, pg.dns)
	require.NoError(t, err)

	ident := pgx.Identifier{t.Name()}.Sanitize()

	defer admin.Close(ctx)

	_, err = admin.Exec(ctx, "CREATE DATABASE "+ident)
	var pgErr *pgconn.PgError
	if err != nil && !(errors.As(err, &pgErr) && pgErr.Code == "42P04") { // duplicate_database
		t.Fatalf("checking database: %v", err)
	}

	u, _ := url.Parse(pg.dns)
	u.Path = "/" + t.Name()
	p, err := pgxpool.New(ctx, u.String())
	require.NoError(t, err)

	return p
}
