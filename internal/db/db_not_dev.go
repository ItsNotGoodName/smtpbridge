//go:build !dev

package db

import (
	"github.com/uptrace/bun"
)

func New(dbPath string) (*bun.DB, error) {
	return new(dbPath)
}
