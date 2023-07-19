package core

import (
	"context"

	"github.com/uptrace/bun"
)

type Context struct {
	Bus  *Bus
	DB   *bun.DB
	File FileStore
	ctx  context.Context
}

func (a App) Context(ctx context.Context) *Context {
	return &Context{
		Bus:  a.Bus,
		DB:   a.DB,
		File: a.File,
		ctx:  ctx,
	}
}

func (c *Context) Context() context.Context {
	return c.ctx
}
