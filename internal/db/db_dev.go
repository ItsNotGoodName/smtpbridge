//go:build dev

package db

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

func New(dbPath string) (*bun.DB, error) {
	db, err := new(dbPath)
	if err != nil {
		return nil, err
	}

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db, nil
}
