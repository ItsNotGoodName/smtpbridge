package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/db/migrations"
	"github.com/uptrace/bun/migrate"
)

func main() {
	migrator := migrate.NewMigrator(nil, migrations.Migrations)

	ctx := context.Background()
	name := strings.Join(os.Args[1:], "_")
	files, err := migrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		panic(err)
	}

	for _, mf := range files {
		fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	}
}
