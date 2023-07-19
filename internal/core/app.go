package core

import (
	"github.com/uptrace/bun"
)

type App struct {
	Bus  *Bus
	DB   *bun.DB
	File FileStore
}

func NewApp(bunDB *bun.DB, fileStore FileStore) App {
	return App{
		Bus:  NewBus(),
		DB:   bunDB,
		File: fileStore,
	}
}
