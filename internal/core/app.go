package core

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/uptrace/bun"
)

type App struct {
	Bus    *Bus
	DB     *bun.DB
	File   FileStore
	Config *models.Config
}

func NewApp(config *models.Config, bunDB *bun.DB, fileStore FileStore) App {
	return App{
		Config: config,
		Bus:    NewBus(),
		DB:     bunDB,
		File:   fileStore,
	}
}

func (a App) newContext(ctx context.Context, actor Actor) Context {
	return Context{
		Bus:    a.Bus,
		DB:     a.DB,
		File:   a.File,
		Config: a.Config,
		ctx:    ctx,
		Actor:  actor,
	}
}

func (a App) Context(ctx context.Context) Context {
	return a.newContext(ctx, ActorAnonymous)
}

func (a App) SystemContext(ctx context.Context) Context {
	return a.newContext(ctx, ActorSystem)
}
