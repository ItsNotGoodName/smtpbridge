package db

import (
	"context"
	"database/sql"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db/migrations"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/migrate"
	_ "modernc.org/sqlite"
)

func new(dbPath string) (*bun.DB, error) {
	// https://github.com/pocketbase/pocketbase/blob/94a1cc07d5ed26bb9ba277506448b4fa14ef4bd9/core/db_nocgo.go#L14
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"
	sqlite, err := sql.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlite, sqlitedialect.New())

	return db, nil
}

func Migrate(cc context.Context, bunDB *bun.DB) error {
	migrator := migrate.NewMigrator(bunDB, migrations.Migrations)

	// Create migrations table
	if err := migrator.Init(cc); err != nil {
		return err
	}

	// Lock
	if err := migrator.Lock(cc); err != nil {
		return err
	}
	defer migrator.Unlock(cc)

	// Migrate
	group, err := migrator.Migrate(cc)
	if err != nil {
		return err
	}

	if group.IsZero() {
		log.Info().Msg("Database migrations are up to date")
	} else {
		log.Info().Msgf("Migrated to %s", group)
	}

	return nil
}

func Size(cc core.Context) (int64, error) {
	var size int64
	err := cc.DB.QueryRow("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();").Scan(&size)
	return size, err
}

func Vacuum(cc core.Context) error {
	_, err := cc.DB.Exec("VACUUM;")
	return err
}
