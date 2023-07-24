package core

import (
	"context"

	"github.com/uptrace/bun"
)

type Context struct {
	Actor Actor
	Bus   *Bus
	DB    *bun.DB
	File  FileStore
	ctx   context.Context
}

func (a App) Context(ctx context.Context) Context {
	return Context{
		Actor: ActorAnon,
		Bus:   a.Bus,
		DB:    a.DB,
		File:  a.File,
		ctx:   ctx,
	}
}

func (c Context) WithActor(actor Actor) Context {
	c.Actor = actor
	return c
}

func (c Context) Context() context.Context {
	return c.ctx
}

type Actor int

const (
	ActorAnon Actor = iota
	ActorUser
	ActorSystem
)
