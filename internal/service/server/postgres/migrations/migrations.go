package migrations

import (
	"embed"
	"fmt"
	"io"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

//go:embed *.sql
var fs embed.FS

func Migrate(dsn string) error {
	pgDSN, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}

	mate := dbmate.New(pgDSN)
	mate.FS = fs
	mate.Log = io.Discard
	mate.MigrationsDir = []string{"./"}

	if err := mate.Wait(); err != nil {
		return fmt.Errorf("wait migrations: %w", err)
	}

	if _, err := mate.FindMigrations(); err != nil {
		return fmt.Errorf("find migrations: %w", err)
	}

	if err := mate.Migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
